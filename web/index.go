package web

const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SyncDevelop</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,"SF Pro Text","Helvetica Neue","Microsoft YaHei",sans-serif;background:#0f0f1a;color:#e0e0e0;min-height:100vh}
.header{background:#161625;padding:16px 24px;display:flex;justify-content:space-between;align-items:center;border-bottom:1px solid #2a2a40}
.header h1{font-size:18px;font-weight:600;background:linear-gradient(135deg,#6c5ce7,#a29bfe);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.status{display:flex;align-items:center;gap:8px;font-size:13px;color:#667}
.dot{width:8px;height:8px;border-radius:50%;background:#00b894;animation:pulse 2s infinite}
.dot.off{background:#e17055;animation:none}
@keyframes pulse{0%,100%{opacity:1}50%{opacity:.4}}

.container{max-width:900px;margin:0 auto;padding:20px}
.card{background:#161625;border-radius:14px;padding:20px;margin-bottom:16px;border:1px solid #2a2a40}
.card h2{font-size:13px;color:#667;margin-bottom:12px;text-transform:uppercase;letter-spacing:1.5px;font-weight:500}

.input-row{display:flex;gap:8px}
.input-area{flex:1;background:#0f0f1a;border:1px solid #2a2a40;border-radius:10px;padding:12px 14px;color:#e0e0e0;font-size:14px;resize:vertical;min-height:72px;outline:none;font-family:inherit;line-height:1.6}
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

.drop-zone{border:2px dashed #2a2a40;border-radius:10px;padding:32px;text-align:center;color:#556;transition:all .2s;cursor:pointer}
.drop-zone:hover,.drop-zone.dragover{border-color:#6c5ce7;color:#a29bfe;background:rgba(108,92,231,.05)}
.drop-zone svg{width:32px;height:32px;margin-bottom:8px;opacity:.4}

.list{display:flex;flex-direction:column;gap:10px;max-height:60vh;overflow-y:auto;padding-right:4px}
.list::-webkit-scrollbar{width:4px}
.list::-webkit-scrollbar-thumb{background:#2a2a40;border-radius:2px}

.item{background:#0f0f1a;border-radius:10px;padding:14px;border:1px solid #2a2a40;transition:border-color .2s;position:relative}
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
.file-info{flex:1}
.file-name{font-size:14px;color:#ddd;font-weight:500}
.file-size{font-size:12px;color:#556}

.empty{text-align:center;color:#445;padding:48px 0;font-size:14px}

.toast{position:fixed;bottom:24px;left:50%;transform:translateX(-50%) translateY(80px);background:#00b894;color:#0f0f1a;padding:8px 24px;border-radius:8px;font-size:14px;font-weight:500;transition:transform .3s;z-index:100;pointer-events:none}
.toast.show{transform:translateX(-50%) translateY(0)}

.device-tag{font-size:11px;background:#2a2a40;padding:2px 8px;border-radius:4px;color:#889}

@media(max-width:640px){
  .container{padding:12px}
  .input-row{flex-direction:column}
  .actions{justify-content:center}
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

<div class="container">
  <div class="card">
    <h2>发送文本</h2>
    <div class="input-row">
      <textarea class="input-area" id="textInput" placeholder="粘贴或输入内容，Ctrl+Enter 发送..."></textarea>
    </div>
    <div class="actions">
      <button class="btn" onclick="sendText()">发送</button>
      <button class="btn btn-ghost" onclick="pasteAndSend()">从剪贴板粘贴</button>
    </div>
  </div>

  <div class="card">
    <h2>传输文件</h2>
    <div class="drop-zone" id="dropZone" onclick="document.getElementById('fileInput').click()">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12"/></svg>
      <div>拖拽文件到此处，或点击选择</div>
      <input type="file" id="fileInput" style="display:none" onchange="uploadFiles(this.files)" multiple>
    </div>
  </div>

  <div class="card">
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
      <h2 style="margin:0">同步记录</h2>
      <button class="btn btn-sm btn-danger" onclick="clearAll()" style="display:none" id="clearBtn">清空</button>
    </div>
    <div class="list" id="itemList">
      <div class="empty">暂无记录</div>
    </div>
  </div>
</div>

<div class="toast" id="toast"></div>

<script>
const deviceId = 'dev-' + Math.random().toString(36).substr(2,6);
document.getElementById('deviceId').textContent = deviceId;
let ws;
let items = [];

// WebSocket 连接
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

// 加载全部条目
async function loadItems() {
  const res = await fetch('/api/items');
  items = await res.json();
  renderItems();
}

// 渲染列表
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
    const typeLabel = item.type === 'text' ? '文本' : item.type === 'file' ? '文件' : '图片';
    let body = '';
    if (item.type === 'text') {
      body = '<div class="item-body" onclick="copyText(this)" title="点击复制">' + escHtml(item.content) + '</div>';
    } else if (item.type === 'file') {
      body = '<div class="item-file"><div class="file-icon">📎</div><div class="file-info"><div class="file-name">' + escHtml(item.fileName) + '</div><div class="file-size">' + fmtSize(item.fileSize) + '</div></div><a href="/api/download/' + item.id + '" class="btn btn-sm" download>下载</a></div>';
    }
    return '<div class="item"><div class="item-head"><span class="item-type ' + item.type + '">' + typeLabel + '</span><div class="item-meta"><span>' + (item.from || '') + '</span><span>' + time + '</span></div></div>' + body + '<div class="item-actions"><button class="btn btn-sm btn-ghost" onclick="deleteItem(\'' + item.id + '\')">删除</button></div></div>';
  }).join('');
}

// 发送文本
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

// 粘贴并发送
async function pasteAndSend() {
  try {
    const text = await navigator.clipboard.readText();
    if (text) {
      document.getElementById('textInput').value = text;
      await sendText();
    }
  } catch { toast('无法读取剪贴板，请手动粘贴'); }
}

// 复制文本
function copyText(el) {
  const text = el.innerText;
  navigator.clipboard.writeText(text).then(() => toast('已复制'));
}

// 上传文件
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

// 删除条目
async function deleteItem(id) {
  await fetch('/api/delete/' + id, {method: 'DELETE'});
  items = items.filter(i => i.id !== id);
  renderItems();
}

// 清空
async function clearAll() {
  if (!confirm('确定清空所有记录？')) return;
  await fetch('/api/clear', {method: 'DELETE'});
  items = [];
  renderItems();
  toast('已清空');
}

// 拖拽
const dz = document.getElementById('dropZone');
dz.addEventListener('dragover', e => {e.preventDefault(); dz.classList.add('dragover')});
dz.addEventListener('dragleave', () => dz.classList.remove('dragover'));
dz.addEventListener('drop', e => {
  e.preventDefault(); dz.classList.remove('dragover');
  uploadFiles(e.dataTransfer.files);
});

// Ctrl+Enter
document.getElementById('textInput').addEventListener('keydown', e => {
  if (e.ctrlKey && e.key === 'Enter') sendText();
});

function toast(msg) {
  const t = document.getElementById('toast');
  t.textContent = msg; t.className = 'toast show';
  setTimeout(() => t.className = 'toast', 1800);
}

function escHtml(s) { return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }
function fmtSize(b) { if(b>=1048576) return (b/1048576).toFixed(1)+'MB'; if(b>=1024) return (b/1024).toFixed(1)+'KB'; return b+'B'; }

// 初始化
connectWS();
loadItems();
</script>
</body>
</html>`
