<script>
  /*
  Clonar Mídia - Interface Principal
  Criado por Erasmo Cardoso - Software Engineer | Electronics Technician
  */
  import { onMount } from 'svelte';
  import { GetDisks, StartClone, IsRoot, ScanPartitions, RecoverFiles, RepairFS, ElevatePrivileges } from '../wailsjs/go/main/App.js';
  import { EventsOn, BrowserOpenURL } from '../wailsjs/runtime/runtime.js';
  import { t, locale } from './i18n.js';

  let disks = [];
  let currentView = 'dashboard'; // 'dashboard', 'ver', 'analisar', 'recuperar', 'clonar', 'config', 'sobre'
  let dragSource = null;
  let sourceDisk = null;
  let targetDisk = null;
  let draggingOver = null; // 'source' ou 'target'
  
  let cloning = false;
  let progress = { bytesCopied: 0, totalBytes: 0, percentage: 0, speed: 0 };
  let statusMessage = "";
  let isRoot = true;

  // Elevação
  let showElevateModal = false;
  let rootPassElevate = "";
  let elevateError = "";

  // Recuperação
  let selectedDiskPath = "";
  let recoveryLog = "";

  let settings = {
    bufferSize: 1024 * 1024 * 64, // 64MB default
  };

  onMount(async () => {
    isRoot = await IsRoot();
    await refreshDisks();
    
    // Listeners de Clonagem
    EventsOn("clone_progress", (p) => { progress = p; });
    EventsOn("clone_log", (log) => { console.log(log); });
    EventsOn("clone_complete", (msg) => {
      cloning = false;
      statusMessage = $t('clone.success');
      refreshDisks();
    });
    EventsOn("clone_error", (err) => {
      cloning = false;
      statusMessage = $t('table.health') + ": " + err;
    });

    // Listeners de Recuperação
    EventsOn("recovery_log", (log) => { recoveryLog += log; });
    EventsOn("recovery_result", (res) => { recoveryLog = res; });
    EventsOn("recovery_complete", (msg) => { recoveryLog += "\n\n✅ " + msg; });
    EventsOn("recovery_error", (err) => { recoveryLog += "\n\n❌ Erro: " + err; });
  });

  async function handleScanPartitions() {
    recoveryLog = "Iniciando TestDisk...\n";
    await ScanPartitions(selectedDiskPath);
  }

  async function handleRecoverFiles() {
    recoveryLog = "Iniciando PhotoRec...\nO processo pode demorar dependendo do tamanho do disco.\n";
    // Por padrão salvamos no home do usuário/Recuperado
    await RecoverFiles(selectedDiskPath, "~/Recuperado");
  }

  async function handleRepairFS() {
    recoveryLog = "Iniciando FSCK...\n";
    await RepairFS(selectedDiskPath);
  }

  async function refreshDisks() {
    try {
      disks = await GetDisks();
    } catch (e) {
      console.error(e);
    }
  }

  function handleDragStart(event, disk) {
    dragSource = { disk };
    event.dataTransfer.setData("text/plain", disk.path);
    event.dataTransfer.effectAllowed = "move";
  }

  function handleDragEnter(event, type) {
    event.preventDefault();
    draggingOver = type;
  }

  function handleDragLeave() {
    draggingOver = null;
  }

  function handleDrop(event, targetType) {
    event.preventDefault();
    draggingOver = null;
    if (!dragSource) return;
    
    if (targetType === 'source') {
      sourceDisk = dragSource.disk;
    } else if (targetType === 'target') {
      targetDisk = dragSource.disk;
    }
    dragSource = null;
  }

  let showConfirmModal = false;
  let selectedDisk = null;
  let rootPassword = "";

  async function startCloning() {
    if (!sourceDisk || !targetDisk) return;
    if (sourceDisk.path === targetDisk.path) {
      statusMessage = $t('clone.error_same');
      return;
    }
    showConfirmModal = true;
  }

  async function confirmCloning() {
    if (!rootPassword) {
      statusMessage = $t('modal.error_pass');
      return;
    }
    showConfirmModal = false;
    cloning = true;
    statusMessage = $t('clone.starting');
    await StartClone({
      source: sourceDisk.path,
      destination: targetDisk.path,
      bufferSize: settings.bufferSize
    }, rootPassword);
    rootPassword = ""; // Limpa por segurança
  }

  async function handleElevate() {
    elevateError = "";
    const success = await ElevatePrivileges(rootPassElevate);
    if (success) {
      isRoot = true;
      showElevateModal = false;
      rootPassElevate = "";
      await refreshDisks();
    } else {
      elevateError = $t('modal.error_elevate');
    }
  }

  function formatSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  function formatSpeed(bytesPerSec) {
    if (bytesPerSec === 0) return '0 B/s';
    return formatSize(bytesPerSec) + '/s';
  }

  function openDonate() {
    BrowserOpenURL("https://www.paypal.com/ncp/payment/8V6WQCGN6HDCQ");
  }
</script>

<div class="app-container">
  <!-- Sidebar -->
  <aside class="sidebar glass">
    <div class="logo-area">
      <h1 class="gradient-text">{$t('app.title')}</h1>
      <p class="signature">{$t('app.signature')}</p>
    </div>

    <nav>
      <button class:active={currentView === 'dashboard'} on:click={() => currentView = 'dashboard'}>
        <span class="icon">🏠</span> {$t('nav.dashboard')}
      </button>
      <button class:active={currentView === 'clonar'} on:click={() => currentView = 'clonar'}>
        <span class="icon">🔄</span> {$t('nav.clone_disk')}
      </button>
      <button class:active={currentView === 'clonar_particao'} on:click={() => currentView = 'clonar_particao'}>
        <span class="icon">🧩</span> {$t('nav.clone_partition')}
      </button>
      <button class:active={currentView === 'ver'} on:click={() => currentView = 'ver'}>
        <span class="icon">📊</span> {$t('nav.view_partitions')}
      </button>
      <button class:active={currentView === 'analisar'} on:click={() => currentView = 'analisar'}>
        <span class="icon">🔍</span> {$t('nav.analyze_health')}
      </button>
      <button class:active={currentView === 'recuperar'} on:click={() => currentView = 'recuperar'}>
        <span class="icon">🛠️</span> {$t('nav.recover')}
      </button>
      <div class="spacer"></div>
      
      <button class="btn-refresh" on:click={refreshDisks}>
        <span class="icon">🔄</span> {$t('nav.refresh')}
      </button>

      <button class:active={currentView === 'config'} on:click={() => currentView = 'config'}>
        <span class="icon">⚙️</span> {$t('nav.settings')}
      </button>
      <button class:active={currentView === 'sobre'} on:click={() => currentView = 'sobre'}>
        <span class="icon">ℹ️</span> {$t('nav.about')}
      </button>

      <button class="btn-donate" on:click={openDonate}>
        <span class="icon">❤️</span> {$t('nav.donate')}
      </button>

      <div class="lang-selector-sidebar">
        <span class="icon">🌐</span>
        <select bind:value={$locale} class="glass-select-mini">
          <option value="pt">PT-BR</option>
          <option value="en">EN-US</option>
        </select>
      </div>

      {#if !isRoot}
        <div class="root-warning glass clickable" on:click={() => showElevateModal = true}>
          {$t('nav.read_only')}
        </div>
      {/if}
    </nav>
  </aside>

  <!-- Main Content -->
  <main class="content">
    {#if currentView === 'dashboard'}
      <header>
        <h2>{$t('dashboard.welcome')}</h2>
        <p>{$t('dashboard.overview')}</p>
      </header>

      <div class="stats-grid">
        <div class="stat-card glass">
          <span class="stat-label">{$t('dashboard.disks_detected')}</span>
          <span class="stat-value">{disks.length}</span>
        </div>
        <div class="stat-card glass">
          <span class="stat-label">{$t('dashboard.total_capacity')}</span>
          <span class="stat-value">{formatSize(disks.reduce((acc, d) => acc + d.size, 0))}</span>
        </div>
        <div class="stat-card glass">
          <span class="stat-label">{$t('dashboard.system_status')}</span>
          <span class="stat-value" class:text-success={isRoot} class:text-warning={!isRoot}>
            {isRoot ? $t('dashboard.status_ready') : $t('dashboard.status_limited')}
          </span>
        </div>
      </div>

      <section class="recent-disks">
        <h3>{$t('dashboard.connected_devices')}</h3>
        <div class="disk-table glass">
          <table>
            <thead>
              <tr>
                <th>{$t('table.model')}</th>
                <th>{$t('table.path')}</th>
                <th>{$t('table.size')}</th>
                <th>{$t('table.health')}</th>
              </tr>
            </thead>
            <tbody>
              {#each disks as disk}
                <tr>
                  <td><strong>{disk.model || 'Desconhecido'}</strong></td>
                  <td><code>{disk.path}</code></td>
                  <td>{formatSize(disk.size)}</td>
                  <td><span class="badge success">{$t('table.healthy')}</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </section>
    {:else if currentView === 'clonar'}
      <header>
        <h2>{$t('clone.title')}</h2>
        <p>{$t('clone.subtitle')}</p>
      </header>

      <section class="clone-zones">
        <div 
          class="zone source-zone" 
          class:active={draggingOver === 'source'}
          on:dragover|preventDefault 
          on:dragenter={(e) => handleDragEnter(e, 'source')}
          on:dragleave={handleDragLeave}
          on:drop={(e) => handleDrop(e, 'source')}
        >
          <h3>{$t('clone.source')}</h3>
          {#if sourceDisk}
            <div class="disk-card mini">
              <strong>{sourceDisk.model || sourceDisk.name}</strong>
              <span>{formatSize(sourceDisk.size)}</span>
            </div>
          {:else}
            <div class="placeholder">{$t('clone.drop_source')}</div>
          {/if}
        </div>

        <div class="arrow">➡️</div>

        <div 
          class="zone target-zone" 
          class:active={draggingOver === 'target'}
          on:dragover|preventDefault 
          on:dragenter={(e) => handleDragEnter(e, 'target')}
          on:dragleave={handleDragLeave}
          on:drop={(e) => handleDrop(e, 'target')}
        >
          <h3>{$t('clone.target')}</h3>
          {#if targetDisk}
            <div class="disk-card mini">
              <strong>{targetDisk.model || targetDisk.name}</strong>
              <span>{formatSize(targetDisk.size)}</span>
            </div>
          {:else}
            <div class="placeholder">{$t('clone.drop_target')}</div>
          {/if}
        </div>
      </section>

      {#if cloning}
        <section class="progress-area glass">
          <div class="progress-header">
            <span>{$t('clone.progress')}: {progress.percentage.toFixed(2)}%</span>
            <span>{formatSpeed(progress.speed)}</span>
          </div>
          <div class="progress-bar-bg">
            <div class="progress-bar-fill" style="width: {progress.percentage}%"></div>
          </div>
          <div class="progress-footer">
            <span>{formatSize(progress.bytesCopied)} de {formatSize(progress.totalBytes)}</span>
            <p class="status">{statusMessage}</p>
          </div>
        </section>
      {:else}
        <div class="actions">
          <button class="btn-primary" disabled={!sourceDisk || !targetDisk} on:click={startCloning}>
            {$t('clone.start')}
          </button>
          <p class="status">{statusMessage}</p>
        </div>
      {/if}

      <section class="disk-list">
        <h3>{$t('clone.available')}</h3>
        <div class="grid">
          {#each disks as disk}
            <div 
              class="disk-card glass" 
              draggable="true" 
              on:dragstart={(e) => handleDragStart(e, disk)}
            >
              <div class="disk-icon">💾</div>
              <div class="disk-info">
                <strong>{disk.model || disk.name}</strong>
                <code>{disk.path}</code>
                <span>{formatSize(disk.size)}</span>
              </div>
            </div>
          {/each}
        </div>
      </section>
    {:else if currentView === 'clonar_particao'}
      <header>
        <h2>{$t('partition.title')}</h2>
        <p>{$t('partition.subtitle')}</p>
      </header>

      <section class="clone-zones">
        <div 
          class="zone source-zone" 
          class:active={draggingOver === 'source'}
          on:dragover|preventDefault 
          on:dragenter={(e) => handleDragEnter(e, 'source')}
          on:dragleave={handleDragLeave}
          on:drop={(e) => handleDrop(e, 'source')}
        >
          <h3>{$t('clone.source')} ({$t('nav.clone_partition')})</h3>
          {#if sourceDisk}
            <div class="disk-card mini">
              <strong>{sourceDisk.name}</strong>
              <span>{sourceDisk.fstype || 'raw'} - {formatSize(sourceDisk.size)}</span>
            </div>
          {:else}
            <div class="placeholder">{$t('partition.drop_source')}</div>
          {/if}
        </div>

        <div class="arrow">➡️</div>

        <div 
          class="zone target-zone" 
          class:active={draggingOver === 'target'}
          on:dragover|preventDefault 
          on:dragenter={(e) => handleDragEnter(e, 'target')}
          on:dragleave={handleDragLeave}
          on:drop={(e) => handleDrop(e, 'target')}
        >
          <h3>{$t('clone.target')} ({$t('nav.clone_partition')})</h3>
          {#if targetDisk}
            <div class="disk-card mini">
              <strong>{targetDisk.name}</strong>
              <span>{targetDisk.fstype || 'raw'} - {formatSize(targetDisk.size)}</span>
            </div>
          {:else}
            <div class="placeholder">{$t('partition.drop_target')}</div>
          {/if}
        </div>
      </section>

      {#if cloning}
        <!-- Mesma área de progresso -->
        <section class="progress-area glass">
          <div class="progress-header">
            <span>{$t('clone.progress')}: {progress.percentage.toFixed(2)}%</span>
            <span>{formatSpeed(progress.speed)}</span>
          </div>
          <div class="progress-bar-bg">
            <div class="progress-bar-fill" style="width: {progress.percentage}%"></div>
          </div>
          <div class="progress-footer">
            <span>{formatSize(progress.bytesCopied)} de {formatSize(progress.totalBytes)}</span>
            <p class="status">{statusMessage}</p>
          </div>
        </section>
      {:else}
        <div class="actions">
          <button class="btn-primary" disabled={!sourceDisk || !targetDisk} on:click={startCloning}>
            {$t('partition.start')}
          </button>
          <p class="status">{statusMessage}</p>
        </div>
      {/if}

      <section class="disk-list">
        <h3>{$t('partition.available')}</h3>
        <div class="grid">
          {#each disks as disk}
            {#each disk.partitions || [] as part}
              <div 
                class="disk-card glass" 
                draggable="true" 
                on:dragstart={(e) => handleDragStart(e, part)}
              >
                <div class="disk-icon">🧩</div>
                <div class="disk-info">
                  <strong>{part.name}</strong>
                  <code>{part.path}</code>
                  <span>{part.fstype || 'Desconhecido'} - {formatSize(part.size)}</span>
                </div>
              </div>
            {/each}
          {/each}
        </div>
      </section>

    {:else if currentView === 'ver'}
      <header>
        <h2>{$t('view.title')}</h2>
        <p>{$t('view.subtitle')}</p>
      </header>

      <div class="partition-layout">
        <aside class="disk-selector glass">
          {#each disks as disk}
            <button class:active={selectedDisk?.path === disk.path} on:click={() => selectedDisk = disk}>
              {disk.model || disk.name}
            </button>
          {/each}
        </aside>

        <section class="partition-details">
          {#if selectedDisk}
            <div class="visual-disk glass">
              <div class="disk-bar">
                {#each selectedDisk.partitions || [] as part}
                  <div 
                    class="part-block" 
                    style="width: {(part.size / selectedDisk.size) * 100}%"
                    title="{part.name} - {formatSize(part.size)}"
                  >
                    {part.fstype || 'raw'}
                  </div>
                {/each}
              </div>
              <div class="part-legend">
                {#each selectedDisk.partitions || [] as part}
                  <div class="legend-item">
                    <span class="dot"></span>
                    <strong>{part.name}</strong>: {formatSize(part.size)} ({part.fstype || 'N/A'})
                  </div>
                {/each}
              </div>
            </div>
          {:else}
            <div class="placeholder">{$t('view.select_disk')}</div>
          {/if}
        </section>
      </div>

    {:else if currentView === 'analisar'}
      <header>
        <h2>{$t('health.title')}</h2>
        <p>{$t('health.subtitle')}</p>
      </header>

      <div class="health-grid">
        {#each disks as disk}
          <div class="health-card glass">
            <div class="health-header">
              <strong>{disk.model || disk.name}</strong>
              <span class="badge success">OK</span>
            </div>
            <div class="health-stats">
              <div class="h-stat"><span>{$t('health.temp')}:</span> <strong>32°C</strong></div>
              <div class="h-stat"><span>{$t('health.hours')}:</span> <strong>1,240h</strong></div>
              <div class="h-stat"><span>{$t('health.reallocated')}:</span> <strong>0</strong></div>
            </div>
          </div>
        {/each}
      </div>

    {:else if currentView === 'recuperar'}
      <header>
        <h2>{$t('recover.title')}</h2>
        <p>{$t('recover.subtitle')}</p>
      </header>

      <div class="recovery-container">
        <section class="recovery-options">
          <div class="tool-card glass">
            <div class="tool-info">
              <h3>🔍 TestDisk</h3>
              <p>{$t('recover.testdisk_desc')}</p>
            </div>
            <div class="tool-actions">
              <select bind:value={selectedDiskPath}>
                <option value="">{$t('recover.select_disk')}</option>
                {#each disks as disk}
                  <option value={disk.path}>{disk.model || disk.name} ({disk.path})</option>
                {/each}
              </select>
              <button class="btn-primary" on:click={handleScanPartitions} disabled={!selectedDiskPath}>
                {$t('recover.scan')}
              </button>
            </div>
          </div>

          <div class="tool-card glass">
            <div class="tool-info">
              <h3>📷 PhotoRec</h3>
              <p>{$t('recover.photorec_desc')}</p>
            </div>
            <div class="tool-actions">
              <button class="btn-primary" on:click={handleRecoverFiles} disabled={!selectedDiskPath}>
                {$t('recover.extract')}
              </button>
            </div>
          </div>

          <div class="tool-card glass">
            <div class="tool-info">
              <h3>🛠️ FSCK</h3>
              <p>{$t('recover.fsck_desc')}</p>
            </div>
            <div class="tool-actions">
              <button class="btn-warning" on:click={handleRepairFS} disabled={!selectedDiskPath}>
                {$t('recover.repair')}
              </button>
            </div>
          </div>
        </section>

        <section class="recovery-console glass">
          <div class="console-header">
            <span>{$t('recover.console')}</span>
            <button class="btn-clear" on:click={() => recoveryLog = ""}>{$t('recover.clear')}</button>
          </div>
          <pre class="console-output">{recoveryLog || $t('recover.waiting')}</pre>
        </section>
      </div>

    {:else if currentView === 'config'}
        <section class="settings glass">
            <h2>{$t('settings.speed')}</h2>
            <div class="field">
                <label>{$t('settings.buffer_size')}</label>
                <select bind:value={settings.bufferSize}>
                    <option value={1024 * 1024 * 10}>{$t('settings.pendrive')}</option>
                    <option value={1024 * 1024 * 64}>{$t('settings.hdd')}</option>
                    <option value={1024 * 1024 * 512}>{$t('settings.ssd')}</option>
                </select>
            </div>
        </section>
    {:else if currentView === 'sobre'}
      <header>
        <h2>{$t('about.title')}</h2>
        <p>{$t('about.subtitle')}</p>
      </header>

      <section class="about-container glass">
        <div class="about-content">
          <h3>{$t('app.title')}</h3>
          <p>{$t('about.desc')}</p>
          
          <div class="features-list">
            <h4>{$t('about.features')}</h4>
            <ul>
              <li>{$t('about.feature1')}</li>
              <li>{$t('about.feature2')}</li>
              <li>{$t('about.feature3')}</li>
              <li>{$t('about.feature4')}</li>
            </ul>
          </div>

          <div class="author-section">
            <h4>{$t('about.author')}</h4>
            <p class="signature-big">{$t('app.signature')}</p>
            <p class="role-desc">Software Engineer | Electronics Technician</p>
          </div>
        </div>
      </section>
    {:else}
      <div class="coming-soon">
        <h2>Funcionalidade em desenvolvimento</h2>
        <p>A aba {currentView} será implementada em breve.</p>
      </div>
    {/if}
  </main>
</div>

{#if showConfirmModal}
  <div class="modal-overlay">
    <div class="modal-content glass">
      <h3>{$t('modal.confirm_title')}</h3>
      <p>{$t('modal.confirm_msg')}</p>
      <div class="confirm-details">
        <div class="conf-item"><strong>{$t('modal.from')}</strong> {sourceDisk.model || sourceDisk.name} ({sourceDisk.path})</div>
        <div class="conf-item"><strong>{$t('modal.to')}</strong> {targetDisk.model || targetDisk.name} ({targetDisk.path})</div>
      </div>
      
      <div class="password-field">
        <label for="root-pass">{$t('modal.root_pass')}</label>
        <input 
          id="root-pass"
          type="password" 
          bind:value={rootPassword} 
          placeholder={$t('modal.root_placeholder')}
          class="glass-input"
        />
      </div>

      <p class="danger-text">{$t('modal.danger')}</p>
      <div class="modal-actions">
        <button class="btn-secondary" on:click={() => showConfirmModal = false}>{$t('modal.cancel')}</button>
        <button class="btn-danger" on:click={confirmCloning}>{$t('modal.confirm_btn')}</button>
      </div>
    </div>
  </div>
{/if}

{#if showElevateModal}
  <div class="modal-overlay">
    <div class="modal-content glass">
      <h3>{$t('modal.elevate_title')}</h3>
      <p>{$t('modal.elevate_msg')}</p>
      
      <div class="password-field">
        <label for="elevate-pass">{$t('modal.root_pass')}</label>
        <input 
          id="elevate-pass"
          type="password" 
          bind:value={rootPassElevate} 
          placeholder={$t('modal.root_placeholder')}
          class="glass-input"
        />
        {#if elevateError}
          <p class="error-text">{elevateError}</p>
        {/if}
      </div>

      <div class="modal-actions">
        <button class="btn-secondary" on:click={() => { showElevateModal = false; elevateError = ""; }}>{$t('modal.cancel')}</button>
        <button class="btn-primary" on:click={handleElevate}>{$t('modal.elevate_btn')}</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .app-container {
    display: flex;
    height: 100vh;
    width: 100vw;
  }

  .sidebar {
    width: 260px;
    padding: 2rem;
    display: flex;
    flex-direction: column;
    border-right: 1px solid rgba(255, 255, 255, 0.1);
  }

  .logo-area h1 {
    font-size: 1.5rem;
    margin-bottom: 0.2rem;
  }

  .signature {
    font-size: 0.7rem;
    opacity: 0.6;
    margin-bottom: 2rem;
  }

  nav {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex-grow: 1;
  }

  nav button {
    background: transparent;
    border: none;
    color: var(--text-color);
    padding: 0.8rem 1rem;
    text-align: left;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 0.8rem;
  }

  nav button:hover, nav button.active {
    background: rgba(56, 189, 248, 0.1);
    color: var(--accent-color);
  }

  .btn-donate {
    margin-top: 0.5rem;
    color: #ff4d4d !important;
    font-weight: bold;
    border: 1px solid rgba(255, 77, 77, 0.2) !important;
  }

  .btn-donate:hover {
    background: rgba(255, 77, 77, 0.1) !important;
    border-color: #ff4d4d !important;
  }

  .spacer { flex-grow: 1; }

  .lang-selector-sidebar {
    display: flex;
    align-items: center;
    gap: 0.8rem;
    padding: 0.8rem 1rem;
    margin-top: 0.5rem;
  }

  .glass-select-mini {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: var(--text-color);
    border-radius: 4px;
    padding: 2px 5px;
    font-size: 0.8rem;
    cursor: pointer;
    outline: none;
  }

  .glass-select-mini option {
    background: #1a1a2e; /* Cor escura para combinar com o tema */
    color: white;
  }

  .content {
    flex-grow: 1;
    padding: 2rem;
    overflow-y: auto;
  }

  header { margin-bottom: 2rem; }

  .clone-zones {
    display: flex;
    align-items: center;
    gap: 2rem;
    margin-bottom: 2rem;
  }

  .zone {
    flex: 1;
    height: 120px;
    border: 2px dashed rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1rem;
    background: rgba(255, 255, 255, 0.02);
  }

  .arrow { font-size: 2rem; opacity: 0.3; }

  .disk-list .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
  }

  .disk-card {
    padding: 1rem;
    border-radius: 12px;
    display: flex;
    gap: 1rem;
    align-items: center;
    cursor: grab;
    transition: transform 0.2s;
  }

  .disk-card:active { cursor: grabbing; }
  .disk-card:hover { transform: translateY(-2px); }

  .disk-icon { font-size: 2rem; }
  .disk-info { display: flex; flex-direction: column; gap: 0.2rem; }
  .disk-info code { font-size: 0.7rem; color: var(--accent-color); }

  .mini {
    background: var(--accent-color);
    color: var(--bg-color);
    width: 100%;
    padding: 0.5rem;
    border-radius: 6px;
    text-align: center;
  }

  .btn-primary {
    background: var(--accent-color);
    color: var(--bg-color);
    border: none;
    padding: 1rem 2rem;
    border-radius: 8px;
    font-weight: bold;
    cursor: pointer;
    font-size: 1rem;
  }

  .btn-primary:disabled { opacity: 0.3; cursor: not-allowed; }

  .progress-area {
    padding: 1.5rem;
    border-radius: 12px;
    margin-bottom: 2rem;
  }

  .progress-bar-bg {
    height: 12px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    margin: 1rem 0;
    overflow: hidden;
  }

  .progress-bar-fill {
    height: 100%;
    background: var(--accent-color);
    transition: width 0.3s;
  }

  .status { margin-top: 1rem; font-size: 0.9rem; font-weight: bold; }

  .settings { padding: 2rem; border-radius: 12px; }
  .field { margin-top: 1.5rem; display: flex; flex-direction: column; gap: 0.5rem; }
  .zone.active {
    border-color: var(--accent-color);
    background: rgba(56, 189, 248, 0.1);
  }
</style>
