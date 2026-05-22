export const PromptHtml = (t) => `
<div id="page-prompt" class="page">
  <div class="page-header">
    <h1 class="page-title">${t('prompt_title')}</h1>
    <p class="page-subtitle">${t('prompt_sub')}</p>
  </div>

  <div class="prompt-area" style="display:flex; flex-direction:column; min-height:400px;">
    <div class="prompt-toolbar" style="display:flex; justify-content:space-between; align-items:center;">
      <div style="display:flex; gap:5px; align-items:center;">
        <button id="tab-normal-prompt" class="btn btn-primary" style="padding:2px 8px; font-size:11px; height:auto; min-height:20px; border-radius:4px;">Normal</button>
        <button id="tab-json-prompt" class="btn btn-outline" style="padding:2px 8px; font-size:11px; height:auto; min-height:20px; border-radius:4px;">JSON</button>
      </div>
      
      <div class="prompt-actions" style="display:flex; align-items:center;">
        <span class="prompt-meta" style="margin-right:15px; font-size:12px;">
          <span class="char-count" id="char-count">0</span> ${t('chars')}
          <span class="char-count" id="json-char-count" style="display:none;">0</span>
        </span>
        <span class="save-status" id="prompt-saved">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" style="width:16px;height:16px;display:inline;vertical-align:text-bottom"><path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" /></svg>
          ${t('saved')}
        </span>
        <button class="btn btn-primary" id="btn-save-prompt">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" /></svg>
          ${t('btn_save_prompt')}
        </button>
        <button class="btn btn-primary" id="btn-save-json-prompt" style="display:none;">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" /></svg>
          Kaydet
        </button>
      </div>
    </div>
    
    <textarea
      id="prompt-editor"
      class="prompt-editor"
      placeholder="${t('prompt_placeholder')}"
      spellcheck="false"
      style="flex:1; width:100%; min-height:300px; resize:vertical;"
    ></textarea>
    
    <textarea
      id="json-prompt-editor"
      class="prompt-editor"
      placeholder="{\n  &quot;role&quot;: &quot;...&quot;\n}"
      spellcheck="false"
      style="display:none; flex:1; width:100%; min-height:300px; resize:vertical; font-family: monospace; font-size: 13px;"
    ></textarea>
  </div>

  <div class="test-bot-area" style="margin-top:20px; padding:20px; background:var(--surface); border:1px solid var(--border); border-radius:12px;">
    <h3 style="font-size:16px; margin:0 0 10px 0;">${t('test_bot_title') || 'Lokal Test'}</h3>
    <p style="font-size:13px; color:var(--text-sec); margin:0 0 15px 0;">${t('test_bot_sub') || 'Promptunuzu kaydettikten sonra lokal olarak deneyin.'}</p>
    <div style="display:flex; gap:10px;">
      <input type="text" id="prompt-test-input" class="form-control" style="flex:1" placeholder="${t('type_message') || 'Test mesajı yazın...'}">
      <button class="btn btn-primary" id="btn-prompt-test">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" /></svg>
        ${t('send') || 'Gönder'}
      </button>
    </div>
  </div>
</div>
`;

