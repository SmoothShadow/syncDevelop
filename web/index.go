package web

const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SyncDevelop</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,"SF Pro Text","Helvetica Neue","Microsoft YaHei",sans-serif;background:#0f0f1a;color:#e0e0e0;min-height:100vh;display:flex;flex-direction:column}

.header{background:#161625;padding:14px 24px;display:flex;justify-content:space-between;align-items:center;border-bottom:1px solid #2a2a40;flex-shrink:0}
.header h1{font-size:18px;font-weight:600;background:linear-gradient(135deg,#6c5ce7,#a29bfe);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.status{display:flex;align-items:center;gap:8px;font-size:13px;color:#667}
.dot{width:8px;height:8px;border-radius:50%;background:#00b894;animation:pulse 2s infinite}
.dot.off{background:#e17055;animation:none}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}

.main{display:flex;flex:1;overflow:hidden}

/* 左侧：同步记录 */
.panel-left{flex:1;display:flex;flex-direction:column;min-width:0;border-right:1px solid #2a2a40}
.panel-head{padding:14px 20px;display:flex;justify-content:space-between;align-items:center;border-bottom:1px solid #2a2a40;flex-shrink:0}
.panel-head h2{font-size:13px;color:#667;text-transform:uppercase;letter-spacing:1.5px;font-weight:500}
.list{flex:1;overflow-y:auto;padding:12px 16px;display:flex;flex-direction:column;gap:10px}
.list::-webkit-scrollbar{width:4px}
.list::-webkit-scrollbar-thumb{background:#2a2a40;border-radius:2px}

/* 右侧：输入区 */
.panel-right{width:380px;flex-shrink:0;display:flex;flex-direction:column;overflow-y:auto}
.panel-right .card{padding:16px 20px;border-bottom:1px solid #2a2a40;border-radius:0;margin:0}
.panel-right .card:last-child{border-bottom:none}

.card h2{font-size:13px;color:#667;margin-bottom:10px;text-transform:uppercase;letter-spacing:1.5px;font-weight:500}

.input-area{width:100%;background:#0f0f1a;border:1px solid #2a2a40;border-radius:10px;padding:12px 14px;color:#e0e0e0;font-size:14px;resize:vertical;min-height:100px;outline:none;font-family:inherit;line-height:1.6}
.input-area:focus{border-color:#6c5ce7}
.input-area::placeholder{color:#445}
.btn{background:linear-gradient(135deg,#6c5ce7,#a29bfe);color:#fff;border:none;padding:10px 24px;border-radius:8px;cursor:pointer;font-size:14px;font-weight:500;transition:opacity .2s;white-space:nowrap}
.btn:hover{opacity:.85}
.btn:active{opacity:.7}
.btn-sm{padding:6px 14px;font-size:12px;border-radius:6px}
.btn-danger{background:linear-gradient(135deg,#e17055,#d63031)}
.btn-ghost{background:transparent;border:1px solid #2a2a40;color:#889}
.btn-ghost:hover{border-color:#6c5ce7;color:#a29bfe}

.actions{display:flex;gap:8px;margin-top:10px;flex-wrap:wrap}

.paste-zone{border:2px dashed #2a2a40;border-radius:10px;padding:20px;text-align:center;color:#556;transition:all .2s;cursor:pointer;position:relative}
.paste-zone:hover,.paste-zone.dragover{border-color:#6c5ce7;color:#a29bfe;background:rgba(108,92,231,.05)}
.paste-zone svg{width:28px;height:28px;margin-bottom:6px;opacity:.4}
.paste-zone .hint{font-size:13px}
.paste-preview{margin-top:10px;text-align:left}
.paste-preview img{max-width:100%;max-height:140px;border-radius:6px;border:1px solid #2a2a40;display:block}
.paste-preview .send-bar{display:flex;align-items:center;gap:8px;margin-top:8px}
.paste-preview .send-bar span{font-size:13px;color:#889;flex:1;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}

/* 条目样式 */
.item{background:#161625;border-radius:10px;padding:14px;border:1px solid #2a2a40;transition:border-color .2s}
.item:hover{border-color:#3a3a55}
.item-head{display:flex;justify-content:space-between;align-items:center;margin-bottom:8px}
.item-type{font-size:11px;padding:3px 10px;border-radius:20px;font-weight:500}
.item-type.text{background:rgba(0,184,148,.12);color:#00b894}
.item-type.file{background:rgba(253,203,110,.12);color:#fdcb6e}
.item-type.image{background:rgba(108,92,231,.12);color:#a29bfe}
.item-meta{font-size:12px;color:#556;display:flex;gap:8px;align-items:center}
.item-body{font-size:14px;white-space:pre-wrap;word-break:break-all;max-height:150px;overflow:hidden;line-height:1.6;color:#ccc;cursor:pointer;padding:4px 0}
.item-body:hover{color:#fff}
.item-actions{display:flex;gap:6px;margin-top:8px}
.item-file{display:flex;align-items:center;gap:10px;padding:4px 0}
.file-icon{width:36px;height:36px;border-radius:8px;background:#2a2a40;display:flex;align-items:center;justify-content:center;font-size:14px}
.file-info{flex:1;min-width:0}
.file-name{font-size:14px;color:#ddd;font-weight:500;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.file-size{font-size:12px;color:#556}
.item-image{cursor:pointer;padding:4px 0}
.item-image img{max-width:100%;max-height:300px;border-radius:8px;border:1px solid #2a2a40;display:block}

.empty{text-align:center;color:#445;padding:48px 0;font-size:14px}

.toast{position:fixed;bottom:24px;left:50%;transform:translateX(-50%) translateY(80px);background:#00b894;color:#0f0f1a;padding:8px 24px;border-radius:8px;font-size:14px;font-weight:500;transition:transform .3s;z-index:100;pointer-events:none}
.toast.show{transform:translateX(-50%) translateY(0)}

.device-tag{font-size:11px;background:#2a2a40;padding:2px 8px;border-radius:4px;color:#889}

.img-modal{display:none;position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,.85);z-index:200;justify-content:center;align-items:center;cursor:zoom-out}
.img-modal.show{display:flex}
.img-modal img{max-width:90vw;max-height:90vh;border-radius:8px}

@media(max-width:768px){
  .main{flex-direction:column}
  .panel-left{border-right:none;border-bottom:1px solid #2a2a40;max-height:45vh}
  .panel-right{width:100%;flex-shrink:1}
}
</style>
</head>
<body>

<div class="header">
  <h1>SyncDevelop</h1>
  <div class="status">
    <div class="dot" id="dot"></div>
    <span id="statusText">连接中</span>
    <span class="device-tag" id="deviceId"></span>
  </div>
</div>

<div class="main">
  <!-- 左侧：同步记录 -->
  <div class="panel-left">
    <div class="panel-head">
      <h2>同步记录</h2>
      <button class="btn btn-sm btn-danger" onclick="clearAll()" style="display:none" id="clearBtn">清空</button>
    </div>
    <div class="list" id="itemList">
      <div class="empty">暂无记录</div>
    </div>
  </div>

  <!-- 右侧：输入区 -->
  <div class="panel-right">
    <div class="card">
      <h2>发送文本</h2>
      <textarea class="input-area" id="textInput" placeholder="粘贴或输入内容，Enter 发送..."></textarea>
      <div class="actions">
        <button class="btn" onclick="sendText()">发送</button>
        <button class="btn btn-ghost" onclick="pasteAndSend()">从剪贴板粘贴</button>
      </div>
    </div>

    <div class="card">
      <h2>粘贴截图 / 传输文件</h2>
      <div class="paste-zone" id="pasteZone" tabindex="0">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12"/></svg>
        <div class="hint">Ctrl+V 粘贴截图 · 拖拽文件 · 点击选择</div>
        <input type="file" id="fileInput" style="display:none" onchange="uploadFiles(this.files)" multiple>
        <div class="paste-preview" id="pastePreview"></div>
      </div>
    </div>
  </div>
</div>

<div class="img-modal" id="imgModal" onclick="this.classList.remove('show')">
  <img id="imgModalSrc" src="">
</div>

<div class="toast" id="toast"></div>

<script>
const deviceId = 'dev-' + Math.random().toString(36).substr(2,6);
document.getElementById('deviceId').textContent = deviceId;
let ws;
let items = [];
let pendingImage = null;

function connectWS() {
  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
  ws = new WebSocket(proto + '//' + location.host + '/ws');
  ws.onopen = () => {
    document.getElementById('dot').className = 'dot';
    document.getElementById('statusText').textContent = '已连接';
  };
  ws.onclose = () => {
    document.getElementById('dot').className = 'dot off';
    document.getElementById('statusText').textContent = '已断开';
    setTimeout(connectWS, 2000);
  };
  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    if (msg.event === 'new_item') {
      items.unshift(msg.item);
      renderItems();
    } else if (msg.event === 'update') {
      loadItems();
    }
  };
}

async function loadItems() {
  const res = await fetch('/api/items');
  items = await res.json();
  renderItems();
}

function renderItems() {
  const list = document.getElementById('itemList');
  const clearBtn = document.getElementById('clearBtn');
  if (!items || items.length === 0) {
    list.innerHTML = '<div class="empty">暂无记录</div>';
    clearBtn.style.display = 'none';
    return;
  }
  clearBtn.style.display = '';
  list.innerHTML = items.map(item => {
    const time = new Date(item.timestamp).toLocaleString();
    const typeLabel = item.type === 'text' ? '文本' : item.type === 'image' ? '图片' : '文件';
    let body = '';
    if (item.type === 'text') {
      body = '<div class="item-body" onclick="copyText(this)" title="点击复制">' + escHtml(item.content) + '</div>';
    } else if (item.type === 'image') {
      body = '<div class="item-image" onclick="previewImg(\'/api/download/' + item.id + '\')"><img src="/api/download/' + item.id + '" loading="lazy" alt="截图"></div>';
    } else if (item.type === 'file') {
      body = '<div class="item-file"><div class="file-icon">📎</div><div class="file-info"><div class="file-name">' + escHtml(item.fileName) + '</div><div class="file-size">' + fmtSize(item.fileSize) + '</div></div><a href="/api/download/' + item.id + '" class="btn btn-sm" download>下载</a></div>';
    }
    return '<div class="item"><div class="item-head"><span class="item-type ' + item.type + '">' + typeLabel + '</span><div class="item-meta"><span>' + (item.from || '') + '</span><span>' + time + '</span></div></div>' + body + '<div class="item-actions"><button class="btn btn-sm btn-ghost" onclick="deleteItem(\'' + item.id + '\')">删除</button></div></div>';
  }).join('');
}

async function sendText() {
  const input = document.getElementById('textInput');
  const text = input.value.trim();
  if (!text) return;
  await fetch('/api/send', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({text, from: deviceId})
  });
  input.value = '';
  toast('已发送');
}

async function pasteAndSend() {
  try {
    const text = await navigator.clipboard.readText();
    if (text) {
      document.getElementById('textInput').value = text;
      await sendText();
    }
  } catch { toast('无法读取剪贴板，请手动粘贴'); }
}

function copyText(el) {
  const text = el.innerText;
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(text).then(() => toast('已复制')).catch(() => fallbackCopy(text));
  } else {
    fallbackCopy(text);
  }
}
function fallbackCopy(text) {
  const ta = document.createElement('textarea');
  ta.value = text;
  ta.style.cssText = 'position:fixed;left:-9999px';
  document.body.appendChild(ta);
  ta.select();
  try { document.execCommand('copy'); toast('已复制'); } catch { toast('复制失败'); }
  document.body.removeChild(ta);
}

// 截图粘贴
document.addEventListener('paste', (e) => {
  const files = e.clipboardData && e.clipboardData.files;
  if (files && files.length > 0) {
    for (const file of files) {
      if (file.type.startsWith('image/')) {
        showPastePreview(file);
        e.preventDefault();
        return;
      }
    }
  }
});

function showPastePreview(file) {
  pendingImage = file;
  const preview = document.getElementById('pastePreview');
  const url = URL.createObjectURL(file);
  preview.innerHTML = '<img src="' + url + '"><div class="send-bar"><span>' + (file.name || 'screenshot.png') + ' (' + fmtSize(file.size) + ')</span><button class="btn btn-sm" onclick="sendPendingImage()">发送截图</button><button class="btn btn-sm btn-ghost" onclick="clearPendingImage()">取消</button></div>';
  document.getElementById('pasteZone').classList.add('dragover');
}

async function sendPendingImage() {
  if (!pendingImage) return;
  const fd = new FormData();
  fd.append('file', pendingImage, pendingImage.name || 'screenshot.png');
  fd.append('from', deviceId);
  await fetch('/api/upload', {method: 'POST', body: fd});
  toast('截图已发送');
  clearPendingImage();
}

function clearPendingImage() {
  pendingImage = null;
  document.getElementById('pastePreview').innerHTML = '';
  document.getElementById('pasteZone').classList.remove('dragover');
}

document.getElementById('pasteZone').addEventListener('click', (e) => {
  if (e.target.closest('.send-bar') || e.target.closest('img')) return;
  document.getElementById('fileInput').click();
});

async function uploadFiles(files) {
  for (const file of files) {
    const fd = new FormData();
    fd.append('file', file);
    fd.append('from', deviceId);
    await fetch('/api/upload', {method: 'POST', body: fd});
    toast('已上传: ' + file.name);
  }
  document.getElementById('fileInput').value = '';
}

function previewImg(src) {
  document.getElementById('imgModalSrc').src = src;
  document.getElementById('imgModal').classList.add('show');
}

async function deleteItem(id) {
  await fetch('/api/delete/' + id, {method: 'DELETE'});
  items = items.filter(i => i.id !== id);
  renderItems();
}

async function clearAll() {
  // if (!confirm('确定清空所有记录？')) return;
  await fetch('/api/clear', {method: 'DELETE'});
  items = [];
  renderItems();
  toast('已清空');
}

const pz = document.getElementById('pasteZone');
pz.addEventListener('dragover', e => {e.preventDefault(); pz.classList.add('dragover')});
pz.addEventListener('dragleave', e => {
  if (!pendingImage) pz.classList.remove('dragover');
});
pz.addEventListener('drop', e => {
  e.preventDefault();
  if (e.dataTransfer.files.length > 0) {
    const imgFile = [...e.dataTransfer.files].find(f => f.type.startsWith('image/'));
    if (imgFile) { showPastePreview(imgFile); }
    else { uploadFiles(e.dataTransfer.files); }
  }
});

document.getElementById('textInput').addEventListener('keydown', e => {
  if (e.key === 'Enter') sendText();
});

function toast(msg) {
  const t = document.getElementById('toast');
  t.textContent = msg; t.className = 'toast show';
  setTimeout(() => t.className = 'toast', 1800);
}
function escHtml(s) { return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }
function fmtSize(b) { if(b>=1048576) return (b/1048576).toFixed(1)+'MB'; if(b>=1024) return (b/1024).toFixed(1)+'KB'; return b+'B'; }

connectWS();
loadItems();
</script>
</body>
</html>`
