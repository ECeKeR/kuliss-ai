export const DashboardHtml = (t) => `
<div id="page-dashboard" class="page active">
  <div class="page-header">
    <h1 class="page-title">${t('Dashboard')}</h1>
    <p class="page-subtitle">${t('dashboard_sub')}</p>
  </div>

  <div class="card status-card" id="status-card">
    <div class="status-left">
      <div class="status-icon-wrap" id="status-icon-wrap">
        <svg id="status-svg" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 3v1.5M4.5 8.25H3m18 0h-1.5M4.5 12H3m18 0h-1.5m-15 3.75H3m18 0h-1.5M8.25 19.5V21M12 3v1.5m0 15V21m3.75-18v1.5m0 15V21m-9-1.5h10.5a2.25 2.25 0 002.25-2.25V6.75a2.25 2.25 0 00-2.25-2.25H6.75A2.25 2.25 0 004.5 6.75v10.5a2.25 2.25 0 002.25 2.25zm.75-12h9v9h-9v-9z" />
        </svg>
      </div>
      <div class="status-text">
        <div class="status-label" id="status-label">${t('status_stopped')}</div>
        <div class="status-sub" id="status-sub">${t('status_start_hint')}</div>
      </div>
    </div>
    <div class="btn-group">
      <button class="btn btn-primary" id="btn-start">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347c-.75.412-1.667-.13-1.667-.986V5.653z" /></svg>
        ${t('btn_start')}
      </button>
      <button class="btn btn-danger" id="btn-stop" disabled>
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M5.25 7.5A2.25 2.25 0 017.5 5.25h9a2.25 2.25 0 012.25 2.25v9a2.25 2.25 0 01-2.25 2.25h-9a2.25 2.25 0 01-2.25-2.25v-9z" /></svg>
        ${t('btn_stop')}
      </button>
    </div>
  </div>

  <div id="qr-area" class="qr-area" style="display:none">
    <div class="card qr-card">
      <span class="qr-badge">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M10.5 1.5H8.25A2.25 2.25 0 006 3.75v16.5a2.25 2.25 0 002.25 2.25h7.5A2.25 2.25 0 0018 20.25V3.75a2.25 2.25 0 00-2.25-2.25H13.5m-3 0V3h3V1.5m-3 0h3m-3 18.75h3" /></svg>
        ${t('qr_waiting')}
      </span>
      <div class="qr-title">${t('qr_scan')}</div>
      <div class="qr-steps">
        ${t('qr_steps')}
      </div>
      <div class="qr-frame">
        <div class="qr-spinner" id="qr-spinner">
          <div class="spinner"></div>
          ${t('qr_preparing')}
        </div>
        <img id="qr-img" src="" alt="QR" class="qr-image" style="display:none" />
      </div>
    </div>
  </div>

  <div class="stats-grid" id="stats-grid">
    <div class="card stat-card">
      <div class="stat-icon">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" /></svg>
      </div>
      <div class="stat-value" id="stat-contacts">—</div>
      <div class="stat-label">${t('stat_contacts')}</div>
    </div>
    <div class="card stat-card">
      <div class="stat-icon">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 01.865-.501 48.172 48.172 0 003.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z" /></svg>
      </div>
      <div class="stat-value" id="stat-msgs">—</div>
      <div class="stat-label">${t('stat_msgs')}</div>
    </div>
    <div class="card stat-card">
      <div class="stat-icon" style="color:var(--danger)">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" /></svg>
      </div>
      <div class="stat-value" id="stat-blocked">—</div>
      <div class="stat-label">${t('stat_blocked')}</div>
    </div>
  </div>
</div>
`;
