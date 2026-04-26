export const BlockedHtml = (t) => `
<div id="page-blocked" class="page">
  <div class="page-header">
    <h1 class="page-title">${t('blocked_title')}</h1>
    <p class="page-subtitle">${t('blocked_sub')}</p>
  </div>
  <div id="blocked-list" class="blocked-list">
    <div class="list-loading">${t('loading')}</div>
  </div>
</div>
`;
