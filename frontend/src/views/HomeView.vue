<template>
  <div class="home-view">
    <section class="hero-panel">
      <div>
        <p class="eyebrow">Live telemetry explorer</p>
        <h2>Welcome to your Formula 1 data hub</h2>
        <p>
          Query OpenF1 datasets.
          Start with a driver overview, then drill into drivers and cars.
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
          <h3>Driver overview</h3>
          <p class="muted">Browse drivers by name or number.</p>
        </div>
        <div class="controls">
          <label>
            Search drivers
            <input v-model="driverSearch" type="text" placeholder="e.g. Norris or 4" />
          </label>
          <button @click="loadOverview" :disabled="overviewLoading">
            {{ overviewLoading ? 'Loading...' : 'Refresh drivers' }}
          </button>
        </div>
      </header>

      <p v-if="overviewError" class="error">{{ overviewError }}</p>

      <div v-if="filteredDrivers.length" class="driver-grid">
        <article
          v-for="driver in highlightedDrivers"
          :key="driver.driver_number"
          class="driver-card"
          @click="setPrimaryDriver(driver)"
        >
          <header>
            <img v-if="driver.headshot_url" :src="driver.headshot_url" :alt="driver.full_name" />
            <div>
              <span class="driver-number">#{{ driver.driver_number }}</span>
              <h4>{{ driver.full_name }}</h4>
              <p class="driver-subtitle">
                {{ driver.first_name }} {{ driver.last_name }}
              </p>
            </div>
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
          </dl>
        </article>
      </div>

      <p v-else-if="overviewDrivers.length && !overviewLoading" class="muted">
        No drivers match that search.
      </p>
      <p v-else-if="!overviewLoading" class="muted">
        No drivers yet. Refresh to load the current driver list.
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
              Driver
              <select v-model="speedDriverNumber">
                <option disabled value="">Select a driver</option>
                <option
                  v-for="driver in filteredDrivers"
                  :key="driver.driver_number"
                  :value="driver.driver_number"
                >
                  #{{ driver.driver_number }} · {{ driver.full_name }}
                </option>
              </select>
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
          Choose a driver to populate this board.
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
              Driver
              <select v-model="gearDriverNumber">
                <option disabled value="">Select a driver</option>
                <option
                  v-for="driver in filteredDrivers"
                  :key="driver.driver_number"
                  :value="driver.driver_number"
                >
                  #{{ driver.driver_number }} · {{ driver.full_name }}
                </option>
              </select>
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
import { DEFAULT_SESSION_KEY } from '../config'

const driverSearch = ref('')
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

const filteredDrivers = computed(() => {
  const query = driverSearch.value.trim().toLowerCase()
  if (!query) return overviewDrivers.value
  return overviewDrivers.value.filter((driver) => {
    const numberMatch = String(driver.driver_number).includes(query)
    const full = `${driver.full_name} ${driver.first_name} ${driver.last_name}`.toLowerCase()
    return numberMatch || full.includes(query)
  })
})

const highlightedDrivers = computed(() => filteredDrivers.value.slice(0, 6))

const speedDriverNumber = ref('')
const speedSummary = ref(null)
const speedLoading = ref(false)
const speedError = ref('')

const gearDriverNumber = ref('')
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

function setPrimaryDriver(driver) {
  if (!driver) return
  speedDriverNumber.value = String(driver.driver_number)
  gearDriverNumber.value = String(driver.driver_number)
  loadSpeedSummary()
  loadGearUsageBoard()
}

async function loadOverview() {
  overviewError.value = ''
  overviewLoading.value = true
  overviewDrivers.value = []
  try {
    overviewDrivers.value = await fetchDrivers({ sessionKey: DEFAULT_SESSION_KEY })
    if (!speedDriverNumber.value && overviewDrivers.value.length) {
      speedDriverNumber.value = String(overviewDrivers.value[0].driver_number)
    }
    if (!gearDriverNumber.value && overviewDrivers.value.length) {
      gearDriverNumber.value = String(overviewDrivers.value[0].driver_number)
    }
  } catch (err) {
    console.error(err)
    overviewError.value = err.message || 'Failed to load drivers'
  } finally {
    overviewLoading.value = false
  }
}

async function loadSpeedSummary() {
  speedError.value = ''
  speedSummary.value = null
  try {
    speedLoading.value = true
    const driverNumber = guardNumber(speedDriverNumber.value, 'driver number')
    speedSummary.value = await fetchDriverSpeedSummary({
      sessionKey: DEFAULT_SESSION_KEY,
      driverNumber,
    })
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
    const driverNumber = guardNumber(gearDriverNumber.value, 'driver number')
    gearUsage.value = await fetchGearUsage({
      sessionKey: DEFAULT_SESSION_KEY,
      driverNumber,
    })
  } catch (err) {
    console.error(err)
    gearError.value = err.message || 'Failed to load gear usage'
  } finally {
    gearLoading.value = false
  }
}

onMounted(() => {
  loadOverview().then(() => {
    loadSpeedSummary()
    loadGearUsageBoard()
  })
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
  border: 1px solid var(--border);
  background: linear-gradient(135deg, rgba(225, 6, 0, 0.18), rgba(15, 15, 18, 0.9));
}

.hero-panel h2 {
  margin: 0.4rem 0 0.6rem;
  font-size: clamp(1.8rem, 3vw, 2.4rem);
}

.hero-panel p {
  max-width: 580px;
  color: var(--muted);
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
  color: #fff;
  background: var(--accent);
}

.hero-actions .secondary {
  background: transparent;
  color: var(--accent);
  border: 1px solid rgba(225, 6, 0, 0.5);
}

.hero-stats {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 0.75rem;
}

.hero-stats li {
  background: rgba(10, 10, 12, 0.9);
  border-radius: 0.75rem;
  padding: 0.9rem;
  border: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.hero-stats span {
  color: var(--muted);
  font-size: 0.85rem;
}

.hero-stats strong {
  font-size: 1.5rem;
}

.eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.2em;
  font-size: 0.75rem;
  color: var(--accent);
}

.panel {
  background: rgba(14, 14, 16, 0.96);
  border-radius: 1rem;
  border: 1px solid var(--border);
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.3);
  padding: 1.5rem;
}

.panel.split {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
}

.panel article {
  background: rgba(8, 8, 10, 0.7);
  padding: 1rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(255, 255, 255, 0.06);
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
  color: var(--muted);
}

.controls input {
  background: var(--bg-0);
  border: 1px solid var(--border);
  border-radius: 0.5rem;
  color: var(--text);
  padding: 0.35rem 0.5rem;
  width: 12rem;
}

.controls select {
  background: var(--bg-0);
  border: 1px solid var(--border);
  border-radius: 0.5rem;
  color: var(--text);
  padding: 0.35rem 0.5rem;
  min-width: 12rem;
}

.controls button {
  background: var(--accent);
  border-radius: 9999px;
  border: none;
  padding: 0.45rem 0.95rem;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
}

.controls button:disabled {
  opacity: 0.6;
  cursor: progress;
}

.muted {
  color: var(--muted);
  font-size: 0.9rem;
}

.error {
  color: var(--accent-2);
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
  background: var(--panel-soft);
  border: 1px solid var(--border);
  cursor: pointer;
  transition: transform 0.15s ease, border-color 0.15s ease;
}

.driver-card:hover {
  transform: translateY(-1px);
  border-color: rgba(225, 6, 0, 0.6);
}

.driver-card header {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  margin-bottom: 0.6rem;
}

.driver-card img {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.driver-number {
  color: var(--accent);
  font-weight: 700;
}

.driver-subtitle {
  margin: 0.1rem 0 0;
  color: var(--muted);
  font-size: 0.8rem;
}

.driver-card dl {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem;
  margin: 0;
}

.driver-card dt {
  font-size: 0.75rem;
  color: var(--muted);
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
  background: var(--panel-soft);
  border-radius: 0.75rem;
  border: 1px solid var(--border);
  padding: 0.85rem;
}

.stat-grid span {
  display: block;
  color: var(--muted);
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
  border-bottom: 1px dashed rgba(255, 255, 255, 0.12);
  padding-bottom: 0.25rem;
}

.gear-count {
  color: var(--accent);
  font-weight: 600;
}

@media (max-width: 900px) {
  .hero-panel {
    grid-template-columns: 1fr;
  }
}
</style>
