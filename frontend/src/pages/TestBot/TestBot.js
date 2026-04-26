export const TestBotHtml = (t) => `
<div id="page-testbot" class="page">
  <div class="page-header">
    <h1 class="page-title" style="display: flex; align-items: center; gap: 10px;">
      <button class="btn btn-secondary" style="font-size: 14px; padding: 6px 12px; display: inline-flex; align-items: center; gap: 6px;" onclick="window.navigateTo('prompt')">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" style="width:16px;height:16px;">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18" />
        </svg>
        ${t('btn_back') || 'Geri'}
      </button>
      ${t('test_bot_title') || 'Test Bot'}
    </h1>
    <p class="page-subtitle">${t('test_bot_sub') || 'Lokalde promptunuzu test edin (WhatsApp bağlantısı gerekmez).'}</p>
  </div>

  <div class="conv-layout">
    <div class="contact-list" id="test-contact-list">
      <div class="contact-list-header">${t('contacts') || 'Kişiler'}</div>
      <div class="contact-item active">
        <div class="contact-top">
          <span class="contact-phone">${t('local_test') || 'Yerel Test'}</span>
        </div>
        <div class="contact-last" id="testbot-contact-last">...</div>
      </div>
    </div>
    <div class="conv-view" style="display:flex; flex-direction:column; background:var(--surface);">
      <div class="conv-header">
        <span class="conv-header-name">${t('local_test') || 'Yerel Test'}</span>
        <div><button class="btn btn-danger" style="font-size:12px" onclick="clearTestChat()">${t('clear') || 'Temizle'}</button></div>
      </div>
      <div class="conv-messages" id="testbot-msgs-container" style="flex:1; overflow-y:auto;">
        <div class="conv-empty">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8.625 12a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H8.25m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0H12m4.125 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm0 0h-.375M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 01-2.555-.337A5.972 5.972 0 015.41 20.97a5.969 5.969 0 01-.474-.065 4.48 4.48 0 00.978-2.025c.09-.457-.133-.901-.467-1.226C3.93 16.178 3 14.189 3 12c0-4.556 4.03-8.25 9-8.25s9 3.694 9 8.25z" />
          </svg>
          <span>${t('test_hint') || 'Mesaj yazarak testi başlatın'}</span>
        </div>
      </div>
      <div class="testbot-input-area" style="padding:15px; border-top:1px solid var(--border); display:flex; gap:10px; background:var(--surface);">
        <input type="text" id="testbot-input" class="form-control" style="flex:1" placeholder="${t('type_message') || 'Mesaj yazın...'}">
        <button class="btn btn-primary" id="btn-testbot-send">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" style="width:20px;height:20px;"><path stroke-linecap="round" stroke-linejoin="round" d="M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5" /></svg>
        </button>
      </div>
    </div>
  </div>
</div>
`;
