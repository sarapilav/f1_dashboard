<template>
  <div class="home-view">
    <section class="hero-panel">
      <div>
        <p class="eyebrow">Live telemetry explorer</p>
        <h2>Welcome to your Formula 1 data hub</h2>
        <p>
          Query OpenF1 datasets.
          Start with a session overview, then drill into drivers and cars.
        </p>
        <div class="hero-actions">
          <RouterLink to="/drivers">Browse drivers</RouterLink>
          <RouterLink to="/cars" class="secondary">Inspect car data</RouterLink>
        </div>
      </div>
      <ul class="hero-stats">
        <li>
          <span>Total drivers</span>
          <strong>{{ overviewStats.totalDrivers }}</strong>
        </li>
        <li>
          <span>Teams represented</span>
          <strong>{{ overviewStats.teams }}</strong>
        </li>
        <li>
          <span>Countries</span>
          <strong>{{ overviewStats.countries }}</strong>
        </li>
      </ul>
    </section>

    <section class="panel">
      <header class="panel-header">
        <div>
          <h3>Session overview</h3>
          <p class="muted">Fetch drivers for a specific session key.</p>
        </div>
        <div class="controls">
          <label>
            Session key
            <input v-model="overviewSessionKey" type="number" min="0" />
          </label>
          <button @click="loadOverview" :disabled="overviewLoading">
            {{ overviewLoading ? 'Loading...' : 'Load session' }}
          </button>
        </div>
      </header>

      <p v-if="overviewError" class="error">{{ overviewError }}</p>

      <div v-if="overviewDrivers.length" class="driver-grid">
        <article v-for="driver in highlightedDrivers" :key="driver.driver_number" class="driver-card">
          <header>
            <span class="driver-number">#{{ driver.driver_number }}</span>
            <h4>{{ driver.full_name }}</h4>
          </header>
          <dl>
            <div>
              <dt>Team</dt>
              <dd>{{ driver.team_name || 'N/A' }}</dd>
            </div>
            <div>
              <dt>Country</dt>
              <dd>{{ driver.country_code || '—' }}</dd>
            </div>
            <div>
              <dt>Session</dt>
              <dd>{{ driver.session_key }}</dd>
            </div>
          </dl>
        </article>
      </div>

      <p v-else-if="!overviewLoading" class="muted">
        No drivers yet. Enter a session key to populate the grid.
      </p>
    </section>

    <section class="panel split">
      <article>
        <header class="panel-header">
          <div>
            <h3>Speed board</h3>
            <p class="muted">Pull aggregated speed stats per driver.</p>
          </div>
          <div class="controls">
            <label>
              Session key
              <input v-model="speedSessionKey" type="number" min="0" />
            </label>
            <label>
              Driver #
              <input v-model="speedDriverNumber" type="number" min="1" />
            </label>
            <button @click="loadSpeedSummary" :disabled="speedLoading">
              {{ speedLoading ? 'Loading...' : 'Update' }}
            </button>
          </div>
        </header>

        <p v-if="speedError" class="error">{{ speedError }}</p>

        <div v-if="speedSummary" class="stat-grid">
          <div>
            <span>Samples</span>
            <strong>{{ speedSummary.samples }}</strong>
          </div>
          <div>
            <span>Min speed</span>
            <strong>{{ speedSummary.min.toFixed(1) }} km/h</strong>
          </div>
          <div>
            <span>Max speed</span>
            <strong>{{ speedSummary.max.toFixed(1) }} km/h</strong>
          </div>
          <div>
            <span>Average</span>
            <strong>{{ speedSummary.avg.toFixed(1) }} km/h</strong>
          </div>
        </div>
        <p v-else-if="!speedLoading" class="muted">
          Choose a session and driver to populate this board.
        </p>
      </article>

      <article>
        <header class="panel-header">
          <div>
            <h3>Gear usage spotlight</h3>
            <p class="muted">See how evenly gears are used in a run.</p>
          </div>
          <div class="controls">
            <label>
              Session key
              <input v-model="gearSessionKey" type="number" min="0" />
            </label>
            <label>
              Driver #
              <input v-model="gearDriverNumber" type="number" min="1" />
            </label>
            <button @click="loadGearUsageBoard" :disabled="gearLoading">
              {{ gearLoading ? 'Loading...' : 'Update' }}
            </button>
          </div>
        </header>

        <p v-if="gearError" class="error">{{ gearError }}</p>

        <ul v-if="gearUsage.length" class="gear-list">
          <li v-for="item in sortedGearUsage" :key="item.gear">
            <span>Gear {{ item.gear }}</span>
            <span class="gear-count">{{ item.count }}</span>
          </li>
        </ul>
        <p v-else-if="!gearLoading" class="muted">
          Fetch data to visualize how each gear is used.
        </p>
      </article>
    </section>
  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import { computed, onMounted, ref } from 'vue'
import { fetchDrivers, fetchDriverSpeedSummary, fetchGearUsage } from '../api/openf1'

const overviewSessionKey = ref('9159')
const overviewDrivers = ref([])
const overviewLoading = ref(false)
const overviewError = ref('')

const overviewStats = computed(() => {
  const teamSet = new Set()
  const countrySet = new Set()
  overviewDrivers.value.forEach((d) => {
    if (d.team_name) teamSet.add(d.team_name)
    if (d.country_code) countrySet.add(d.country_code)
  })
  return {
    totalDrivers: overviewDrivers.value.length || '—',
    teams: teamSet.size || '—',
    countries: countrySet.size || '—',
  }
})

const highlightedDrivers = computed(() => overviewDrivers.value.slice(0, 6))

const speedSessionKey = ref('9159')
const speedDriverNumber = ref('55')
const speedSummary = ref(null)
const speedLoading = ref(false)
const speedError = ref('')

const gearSessionKey = ref('9159')
const gearDriverNumber = ref('55')
const gearUsage = ref([])
const gearLoading = ref(false)
const gearError = ref('')

const sortedGearUsage = computed(() => [...gearUsage.value].sort((a, b) => a.gear - b.gear))

function guardNumber(value, label) {
  if (!value && value !== 0) {
    throw new Error(`Please provide a ${label}.`)
  }
  const parsed = Number(value)
  if (Number.isNaN(parsed)) {
    throw new Error(`${label} must be a number.`)
  }
  return parsed
}

async function loadOverview() {
  overviewError.value = ''
  overviewLoading.value = true
  overviewDrivers.value = []
  try {
    const sessionKey = guardNumber(overviewSessionKey.value, 'session key')
    overviewDrivers.value = await fetchDrivers({ sessionKey })
  } catch (err) {
    console.error(err)
    overviewError.value = err.message || 'Failed to load session'
  } finally {
    overviewLoading.value = false
  }
}

async function loadSpeedSummary() {
  speedError.value = ''
  speedSummary.value = null
  try {
    speedLoading.value = true
    const sessionKey = guardNumber(speedSessionKey.value, 'session key')
    const driverNumber = guardNumber(speedDriverNumber.value, 'driver number')
    speedSummary.value = await fetchDriverSpeedSummary({ sessionKey, driverNumber })
  } catch (err) {
    console.error(err)
    speedError.value = err.message || 'Failed to load speed summary'
  } finally {
    speedLoading.value = false
  }
}

async function loadGearUsageBoard() {
  gearError.value = ''
  gearUsage.value = []
  try {
    gearLoading.value = true
    const sessionKey = guardNumber(gearSessionKey.value, 'session key')
    const driverNumber = guardNumber(gearDriverNumber.value, 'driver number')
    gearUsage.value = await fetchGearUsage({ sessionKey, driverNumber })
  } catch (err) {
    console.error(err)
    gearError.value = err.message || 'Failed to load gear usage'
  } finally {
    gearLoading.value = false
  }
}

onMounted(() => {
  loadOverview()
  loadSpeedSummary()
  loadGearUsageBoard()
})
</script>

<style scoped>
.home-view {
  display: flex;
  flex-direction: column;
  gap: 1.75rem;
}

.hero-panel {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
  gap: 1.5rem;
  padding: 1.75rem;
  border-radius: 1rem;
  border: 1px solid #1f2937;
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.15), rgba(79, 70, 229, 0.15));
}

.hero-panel h2 {
  margin: 0.4rem 0 0.6rem;
  font-size: clamp(1.8rem, 3vw, 2.4rem);
}

.hero-panel p {
  max-width: 580px;
  color: #cbd5f5;
}

.hero-actions {
  margin-top: 1rem;
  display: flex;
  gap: 0.5rem;
}

.hero-actions a {
  text-decoration: none;
  padding: 0.55rem 1rem;
  border-radius: 9999px;
  font-weight: 600;
  color: #020617;
  background: #facc15;
}

.hero-actions .secondary {
  background: transparent;
  color: #facc15;
  border: 1px solid rgba(250, 204, 21, 0.5);
}

.hero-stats {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 0.75rem;
}

.hero-stats li {
  background: rgba(2, 6, 23, 0.75);
  border-radius: 0.75rem;
  padding: 0.9rem;
  border: 1px solid rgba(148, 163, 184, 0.2);
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.hero-stats span {
  color: #94a3b8;
  font-size: 0.85rem;
}

.hero-stats strong {
  font-size: 1.5rem;
}

.eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.2em;
  font-size: 0.75rem;
  color: #a5b4fc;
}

.panel {
  background: rgba(15, 23, 42, 0.95);
  border-radius: 1rem;
  border: 1px solid #1f2937;
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.3);
  padding: 1.5rem;
}

.panel.split {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
}

.panel article {
  background: rgba(2, 6, 23, 0.5);
  padding: 1rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(148, 163, 184, 0.15);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
  margin-bottom: 1rem;
}

.panel-header h3 {
  margin: 0;
}

.controls {
  display: flex;
  flex-wrap: wrap;
  gap: 0.6rem;
  font-size: 0.85rem;
  align-items: flex-end;
}

.controls label {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  color: #94a3b8;
}

.controls input {
  background: #020617;
  border: 1px solid #374151;
  border-radius: 0.5rem;
  color: #e5e7eb;
  padding: 0.35rem 0.5rem;
  width: 7rem;
}

.controls button {
  background: #0ea5e9;
  border-radius: 9999px;
  border: none;
  padding: 0.45rem 0.95rem;
  color: #0b1120;
  font-weight: 600;
  cursor: pointer;
}

.controls button:disabled {
  opacity: 0.6;
  cursor: progress;
}

.muted {
  color: #94a3b8;
  font-size: 0.9rem;
}

.error {
  color: #f87171;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}

.driver-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 0.9rem;
}

.driver-card {
  border-radius: 0.85rem;
  padding: 0.9rem;
  background: #020617;
  border: 1px solid #111827;
}

.driver-card header {
  display: flex;
  align-items: baseline;
  gap: 0.4rem;
  margin-bottom: 0.6rem;
}

.driver-number {
  color: #a5b4fc;
  font-weight: 700;
}

.driver-card dl {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem;
  margin: 0;
}

.driver-card dt {
  font-size: 0.75rem;
  color: #94a3b8;
}

.driver-card dd {
  margin: 0;
  font-weight: 500;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.75rem;
}

.stat-grid div {
  background: #020617;
  border-radius: 0.75rem;
  border: 1px solid #111827;
  padding: 0.85rem;
}

.stat-grid span {
  display: block;
  color: #94a3b8;
  font-size: 0.8rem;
}

.stat-grid strong {
  font-size: 1.3rem;
}

.gear-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.gear-list li {
  display: flex;
  justify-content: space-between;
  border-bottom: 1px dashed rgba(148, 163, 184, 0.3);
  padding-bottom: 0.25rem;
}

.gear-count {
  color: #fcd34d;
  font-weight: 600;
}

@media (max-width: 900px) {
  .hero-panel {
    grid-template-columns: 1fr;
  }
}
</style>
