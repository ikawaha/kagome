function cb(data, status) {
    if (status === "success" && Array.isArray(data.tokens)) {
        $("#morphs").empty();
        $.each(data.tokens, function(i, val) {
            var pos     = (val.pos == null || val.pos.length === 0) ? "*" : val.pos.join(",");
            var base    = val.base_form     != "" ? val.base_form     : "*";
            var reading = val.reading       != "" ? val.reading       : "*";
            var pronoun = val.pronunciation != "" ? val.pronunciation : "*";
            var row = $("<tr>").attr("data-morph-idx", i);
            row.append($("<td>").text(val.surface));
            row.append($("<td>").text(pos));
            row.append($("<td>").text(base));
            row.append($("<td>").text(reading));
            row.append($("<td>").text(pronoun));
            row.click(function() {
                var idx = parseInt($(this).attr("data-morph-idx"), 10);
                focusLatticeNode($(this).hasClass('morph-selected') ? -1 : idx);
            });
            $("#morphs").append(row);
        });
    }
}

function focusLatticeNode(idx) {
    var svgEl = document.querySelector('#lattice-output svg');

    // clear table row highlight
    document.querySelectorAll('#morphs tr.morph-selected').forEach(function(tr) {
        tr.classList.remove('morph-selected');
    });
    var rows = document.querySelectorAll('#morphs tr');
    if (idx >= 0 && idx < rows.length) {
        rows[idx].classList.add('morph-selected');
    }

    if (!svgEl || !cachedBestNodes) return;

    // clear previous SVG highlight
    svgEl.querySelectorAll('g.node.morph-highlight').forEach(function(g) {
        g.classList.remove('morph-highlight');
    });

    if (idx < 0 || idx >= cachedBestNodes.length) return;

    var targetNode = cachedBestNodes[idx].el;
    targetNode.classList.add('morph-highlight');

    // scroll the node into view within #lattice-output (horizontal and vertical)
    var ellipse = targetNode.querySelector('ellipse');
    if (ellipse) {
        var outEl = document.getElementById('lattice-output');
        var elRect = ellipse.getBoundingClientRect();
        var outRect = outEl.getBoundingClientRect();
        if (elRect.left < outRect.left || elRect.right > outRect.right) {
            outEl.scrollLeft += elRect.left - outRect.left - (outRect.width - elRect.width) / 2;
        }
        if (elRect.top < outRect.top || elRect.bottom > outRect.bottom) {
            outEl.scrollTop += elRect.top - outRect.top - (outRect.height - elRect.height) / 2;
        }
    }
}

function tokenize() {
    var s = document.getElementById("inp").value;
    var m = $('input[name="r"]').filter(':checked').val();
    var o = {"sentence": s, "mode": m};
    $.post('./tokenize', JSON.stringify(o), cb, 'json');
}

// zoom
var currentZoom = 1.0;
var svgOrigDims = null;
var svgOrigPx   = null;
var cachedBestNodes = null;
var cachedSVGString = null;
var zoomStep = 0.25;
var minZoom = 0.25;
var maxZoom = 4.0;

function toPx(value, unit) {
    return unit === 'pt' ? value * 96 / 72 : value;
}

function scrollToSVGContent() {
    var outEl = document.getElementById('lattice-output');
    var svgEl = outEl && outEl.querySelector('svg');
    if (!svgEl) return;
    var polygon = svgEl.querySelector('polygon');
    if (polygon) {
        var emptyTop = polygon.getBoundingClientRect().top - svgEl.getBoundingClientRect().top;
        if (emptyTop > 0) {
            outEl.scrollTop = emptyTop;
        }
    }
}

function getZoomAnchor(outEl, svgEl) {
    if (!svgOrigPx) return null;
    var curW = svgOrigPx.width  * currentZoom;
    var curH = svgOrigPx.height * currentZoom;
    // 選択中の形態素ノードがあれば、その中心を基準にする
    // getBoundingClientRect() を使うことで graphviz の内部座標系に依存せず正確な位置を取得する
    var highlight = svgEl && svgEl.querySelector('g.node.morph-highlight ellipse');
    if (highlight) {
        var elRect  = highlight.getBoundingClientRect();
        var outRect = outEl.getBoundingClientRect();
        var elCenterX = (elRect.left + elRect.right)  / 2 - outRect.left + outEl.scrollLeft;
        var elCenterY = (elRect.top  + elRect.bottom) / 2 - outRect.top  + outEl.scrollTop;
        return {
            fracX: elCenterX / curW,
            fracY: elCenterY / curH,
            viewportOffsetX: outEl.clientWidth  / 2,
            viewportOffsetY: outEl.clientHeight / 2,
        };
    }
    // なければ表示領域の中心を基準にする
    return {
        fracX: (outEl.scrollLeft + outEl.clientWidth  / 2) / curW,
        fracY: (outEl.scrollTop  + outEl.clientHeight / 2) / curH,
        viewportOffsetX: outEl.clientWidth  / 2,
        viewportOffsetY: outEl.clientHeight / 2,
    };
}

function applyZoom(anchor) {
    var svgEl = document.querySelector('#lattice-output svg');
    var outEl = document.getElementById('lattice-output');
    if (!svgEl || !svgOrigDims) return;
    svgEl.setAttribute('width',  (svgOrigDims.width  * currentZoom) + svgOrigDims.wUnit);
    svgEl.setAttribute('height', (svgOrigDims.height * currentZoom) + svgOrigDims.hUnit);
    document.getElementById('zoom-label').textContent = Math.round(currentZoom * 100) + '%';
    if (anchor && svgOrigPx) {
        var newW = svgOrigPx.width  * currentZoom;
        var newH = svgOrigPx.height * currentZoom;
        outEl.scrollLeft = Math.max(0, newW * anchor.fracX - anchor.viewportOffsetX);
        outEl.scrollTop  = Math.max(0, newH * anchor.fracY - anchor.viewportOffsetY);
    } else {
        scrollToSVGContent();
    }
}

function zoomIn() {
    var outEl = document.getElementById('lattice-output');
    var svgEl = outEl && outEl.querySelector('svg');
    var anchor = getZoomAnchor(outEl, svgEl);
    currentZoom = Math.min(maxZoom, parseFloat((currentZoom + zoomStep).toFixed(2)));
    applyZoom(anchor);
}

function zoomOut() {
    var outEl = document.getElementById('lattice-output');
    var svgEl = outEl && outEl.querySelector('svg');
    var anchor = getZoomAnchor(outEl, svgEl);
    currentZoom = Math.max(minZoom, parseFloat((currentZoom - zoomStep).toFixed(2)));
    applyZoom(anchor);
}


function fitZoom() {
    if (!svgOrigPx) return;
    var outEl = document.getElementById('lattice-output');
    var containerW = outEl.clientWidth;
    var containerH = parseInt(window.getComputedStyle(outEl).maxHeight) || outEl.clientHeight;
    var zoom = Math.min(containerW / svgOrigPx.width, containerH / svgOrigPx.height);
    currentZoom = parseFloat(Math.max(minZoom, Math.min(maxZoom, zoom)).toFixed(2));
    applyZoom();
}

document.getElementById('lattice-output').addEventListener('wheel', function(e) {
    if (e.ctrlKey || e.metaKey) {
        e.preventDefault();
        if (e.deltaY < 0) { zoomIn(); } else { zoomOut(); }
    }
}, { passive: false });

function updateLattice() {
    var s = document.getElementById("inp").value;
    if (s.trim() === '') {
        document.getElementById('lattice-error').textContent = '';
        document.getElementById('lattice-output').innerHTML = '';
        document.getElementById('lattice-controls').style.display = 'none';
        svgOrigDims = null;
        svgOrigPx   = null;
        cachedBestNodes = null;
        cachedSVGString = null;
        return;
    }
    var m = $('input[name="r"]').filter(':checked').val();
    $.post('./lattice', {s: s, r: m}, function(data) {
        var errEl = document.getElementById('lattice-error');
        var outEl = document.getElementById('lattice-output');
        var ctrlEl = document.getElementById('lattice-controls');
        if (data.error) {
            errEl.textContent = 'Error: ' + data.error;
            outEl.innerHTML = '';
            ctrlEl.style.display = 'none';
            svgOrigDims = null;
            svgOrigPx   = null;
            cachedBestNodes = null;
        } else {
            errEl.textContent = '';
            outEl.innerHTML = data.svg || '';
            cachedSVGString = data.svg || null;
            var svgEl = outEl.querySelector('svg');
            if (svgEl) {
                var wAttr = svgEl.getAttribute('width')  || '';
                var hAttr = svgEl.getAttribute('height') || '';
                var wM = wAttr.match(/^([0-9.]+)([a-z]*)$/i);
                var hM = hAttr.match(/^([0-9.]+)([a-z]*)$/i);
                svgOrigDims = {
                    width:  wM ? parseFloat(wM[1]) : 100,
                    height: hM ? parseFloat(hM[1]) : 100,
                    wUnit:  wM ? wM[2] : '',
                    hUnit:  hM ? hM[2] : '',
                };
                svgOrigPx = {
                    width:  toPx(svgOrigDims.width,  svgOrigDims.wUnit),
                    height: toPx(svgOrigDims.height, svgOrigDims.hUnit),
                };
                // best-path nodes have peripheries=2 in DOT → two <ellipse> in SVG
                // BOS/EOS are also in the optimal path and get peripheries=2, so exclude them by label
                cachedBestNodes = [];
                svgEl.querySelectorAll('g.node').forEach(function(g) {
                    var ellipses = g.querySelectorAll('ellipse');
                    if (ellipses.length >= 2) {
                        var firstText = g.querySelector('text');
                        if (!firstText) return;
                        var label = firstText.textContent.trim();
                        if (label === 'BOS' || label === 'EOS') return;
                        cachedBestNodes.push({ el: g, cx: parseFloat(ellipses[0].getAttribute('cx') || '0') });
                    }
                });
                cachedBestNodes.sort(function(a, b) { return a.cx - b.cx; });
                ctrlEl.style.display = '';
                fitZoom();
            }
        }
    }, 'json');
}

function toggleDownloadMenu(event) {
    event.stopPropagation();
    var menu = document.getElementById('download-menu');
    menu.classList.toggle('open');
}

function closeDownloadMenu() {
    var menu = document.getElementById('download-menu');
    if (menu) menu.classList.remove('open');
}

function downloadAs(format) {
    closeDownloadMenu();
    if (!cachedSVGString) return;

    var inp = document.getElementById('inp').value.trim();
    var base = inp
        ? inp.replace(/[^0-9A-Za-z\u3040-\u9FFF]/g, '_').slice(0, 50)
        : 'lattice';

    if (format === 'svg') {
        var blob = new Blob([cachedSVGString], {type: 'image/svg+xml'});
        var url = URL.createObjectURL(blob);
        var a = document.createElement('a');
        a.href = url;
        a.download = base + '.svg';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    } else if (format === 'png') {
        var w = svgOrigPx ? svgOrigPx.width  : 800;
        var h = svgOrigPx ? svgOrigPx.height : 600;
        var scale = 2;
        var blob2 = new Blob([cachedSVGString], {type: 'image/svg+xml'});
        var url2 = URL.createObjectURL(blob2);
        var img = new Image();
        img.onload = function() {
            var canvas = document.createElement('canvas');
            canvas.width  = w * scale;
            canvas.height = h * scale;
            var ctx = canvas.getContext('2d');
            ctx.scale(scale, scale);
            ctx.drawImage(img, 0, 0, w, h);
            URL.revokeObjectURL(url2);
            var a2 = document.createElement('a');
            a2.href = canvas.toDataURL('image/png');
            a2.download = base + '.png';
            document.body.appendChild(a2);
            a2.click();
            document.body.removeChild(a2);
        };
        img.src = url2;
    }
}

var latticeTimer = null;
function scheduleLattice() {
    if (latticeTimer !== null) {
        clearTimeout(latticeTimer);
    }
    latticeTimer = setTimeout(function() {
        latticeTimer = null;
        updateLattice();
    }, 500);
}

$('input[name="r"]:radio').change(function() {
    var s = document.getElementById("inp").value;
    var m = $('input[name="r"]').filter(':checked').val();
    var o = {"sentence": s, "mode": m};
    $.post('./tokenize', JSON.stringify(o), cb, 'json');
    updateLattice();
});

document.addEventListener('click', function() {
    closeDownloadMenu();
});
