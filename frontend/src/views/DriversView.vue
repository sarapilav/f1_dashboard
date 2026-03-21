<template>
  <div class="drivers-view">
    <section class="panel hero">
      <div>
        <h2>Driver intelligence</h2>
        <p class="muted">
          Browse drivers by name or number, then focus on a specific competitor to reveal speed
          summaries, gear usage, and raw telemetry samples.
        </p>
      </div>
      <ul class="tips">
        <li>Type a name or number to filter the driver list.</li>
        <li>You can chain queries quickly from keyboard (Enter must focus button).</li>
        <li>The telemetry preview shows the first five points from /api/car-data/raw.</li>
      </ul>
    </section>

    <section class="panel">
      <div class="panel-header">
        <h2>Drivers</h2>
        <div class="controls">
          <label>
            Search
            <input v-model="driverSearch" type="text" placeholder="e.g. Hamilton or 44" />
          </label>
          <button @click="loadDrivers" :disabled="loadingDrivers">
            {{ loadingDrivers ? 'Loading...' : 'Refresh drivers' }}
          </button>
        </div>
      </div>

      <div v-if="driversError" class="error">
        {{ driversError }}
      </div>

      <div class="drivers-list" v-if="filteredDrivers.length">
        <div
          v-for="driver in filteredDrivers"
          :key="driver.driver_number"
          :class="[
            'driver-card',
            driver.driver_number === selectedDriverNumber ? 'driver-card--selected' : '',
          ]"
          @click="selectDriver(driver)"
        >
          <div class="driver-header">
            <img v-if="driver.headshot_url" :src="driver.headshot_url" :alt="driver.full_name" />
            <div>
              <span class="driver-number">#{{ driver.driver_number }}</span>
              <span class="driver-name">{{ driver.full_name }}</span>
              <span class="driver-subtitle"> {{ driver.first_name }} {{ driver.last_name }} </span>
            </div>
          </div>
          <div class="driver-meta">
            <span class="team-pill" :style="teamColorStyle(driver.team_colour)">
              {{ driver.team_name || 'Unknown team' }}
            </span>
            <span class="country-code">{{ driver.country_code }}</span>
          </div>
        </div>
      </div>

      <div v-else-if="drivers.length && !loadingDrivers">
        <p class="muted">No drivers match that search.</p>
      </div>
      <div v-else-if="!loadingDrivers">
        <p class="muted">No drivers loaded yet. Refresh to load the current driver list.</p>
      </div>
    </section>

    <section class="panel" v-if="selectedDriverNumber">
      <div class="panel-header">
        <h2>Driver details</h2>
        <p class="muted">Driver #{{ selectedDriverNumber }} · {{ selectedDriverName }}</p>
      </div>

      <div v-if="driverSummaryError" class="error">
        {{ driverSummaryError }}
      </div>

      <div class="driver-layout">
        <div class="summary-cards">
          <div class="card driver-profile" v-if="selectedDriver">
            <div class="profile-header">
              <img
                v-if="selectedDriver.headshot_url"
                :src="selectedDriver.headshot_url"
                :alt="selectedDriver.full_name"
              />
              <div>
                <h3>{{ selectedDriver.full_name }}</h3>
                <p class="muted">
                  #{{ selectedDriver.driver_number }} · {{ selectedDriver.team_name || 'Team TBD' }}
                </p>
              </div>
            </div>
            <dl>
              <div>
                <dt>First name</dt>
                <dd>{{ selectedDriver.first_name }}</dd>
              </div>
              <div>
                <dt>Last name</dt>
                <dd>{{ selectedDriver.last_name }}</dd>
              </div>
              <div>
                <dt>Country</dt>
                <dd>{{ selectedDriver.country_code || '—' }}</dd>
              </div>
            </dl>
          </div>

          <div class="card" v-if="speedSummary">
            <h3>Speed summary</h3>
            <div class="card-grid">
              <div class="metric">
                <span class="metric-label">Samples</span>
                <span class="metric-value">{{ speedSummary.samples }}</span>
              </div>
              <div class="metric">
                <span class="metric-label">Min speed</span>
                <span class="metric-value">{{ speedSummary.min.toFixed(1) }} km/h</span>
              </div>
              <div class="metric">
                <span class="metric-label">Max speed</span>
                <span class="metric-value">{{ speedSummary.max.toFixed(1) }} km/h</span>
              </div>
              <div class="metric">
                <span class="metric-label">Avg speed</span>
                <span class="metric-value">{{ speedSummary.avg.toFixed(1) }} km/h</span>
              </div>
            </div>
          </div>

          <div class="card" v-if="gearUsage.length">
            <h3>Gear usage</h3>
            <ul class="gear-list">
              <li v-for="item in sortedGearUsage" :key="item.gear">
                <span>Gear {{ item.gear }}</span>
                <span class="gear-count">{{ item.count }}</span>
              </li>
            </ul>
          </div>
        </div>

        <div class="card">
          <div class="card-header">
            <h3>Telemetry charts</h3>
            <button @click="loadCarData" :disabled="loadingCarData">
              {{ loadingCarData ? 'Loading telemetry...' : 'Fetch raw telemetry' }}
            </button>
          </div>
          <TelemetryCharts :telemetry="carData" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchDrivers, fetchDriverSpeedSummary, fetchGearUsage, fetchCarData } from '../api/openf1'
import { DEFAULT_SESSION_KEY } from '../config'
import TelemetryCharts from '../components/TelemetryCharts.vue'

const driverSearch = ref('')
const drivers = ref([])
const loadingDrivers = ref(false)
const driversError = ref('')

const selectedDriverNumber = ref(null)
const selectedDriverName = ref('')

const speedSummary = ref(null)
const gearUsage = ref([])
const driverSummaryError = ref('')

const carData = ref([])
const loadingCarData = ref(false)

const selectedDriver = computed(() =>
  drivers.value.find((driver) => driver.driver_number === selectedDriverNumber.value),
)

const filteredDrivers = computed(() => {
  const query = driverSearch.value.trim().toLowerCase()
  if (!query) return drivers.value
  return drivers.value.filter((driver) => {
    const numberMatch = String(driver.driver_number).includes(query)
    const full = `${driver.full_name} ${driver.first_name} ${driver.last_name}`.toLowerCase()
    return numberMatch || full.includes(query)
  })
})

const sortedGearUsage = computed(() => [...gearUsage.value].sort((a, b) => a.gear - b.gear))

function teamColorStyle(hex) {
  if (!hex) return {}
  const normalized = hex.startsWith('#') ? hex : `#${hex}`
  return {
    backgroundColor: normalized,
    borderColor: normalized,
  }
}

async function loadDrivers() {
  driversError.value = ''
  loadingDrivers.value = true
  drivers.value = []
  selectedDriverNumber.value = null
  speedSummary.value = null
  gearUsage.value = []
  carData.value = []

  try {
    drivers.value = await fetchDrivers({ sessionKey: DEFAULT_SESSION_KEY })
    if (!selectedDriverNumber.value && drivers.value.length) {
      selectDriver(drivers.value[0])
    }
  } catch (err) {
    console.error(err)
    driversError.value = err.message || 'Failed to load drivers'
  } finally {
    loadingDrivers.value = false
  }
}

async function selectDriver(driver) {
  selectedDriverNumber.value = driver.driver_number
  selectedDriverName.value = driver.full_name
  await Promise.all([loadSpeedSummary(), loadGearUsage()])
}

async function loadSpeedSummary() {
  driverSummaryError.value = ''
  speedSummary.value = null

  try {
    const summary = await fetchDriverSpeedSummary({
      driverNumber: selectedDriverNumber.value,
      sessionKey: DEFAULT_SESSION_KEY,
    })
    speedSummary.value = summary
  } catch (err) {
    console.error(err)
    driverSummaryError.value = err.message || 'Failed to load speed summary'
  }
}

async function loadGearUsage() {
  gearUsage.value = []

  try {
    const data = await fetchGearUsage({
      driverNumber: selectedDriverNumber.value,
      sessionKey: DEFAULT_SESSION_KEY,
    })
    gearUsage.value = data
  } catch (err) {
    console.error(err)
  }
}

async function loadCarData() {
  loadingCarData.value = true
  carData.value = []
  try {
    const data = await fetchCarData({
      driverNumber: selectedDriverNumber.value,
      sessionKey: DEFAULT_SESSION_KEY,
    })
    carData.value = data
  } catch (err) {
    console.error(err)
  } finally {
    loadingCarData.value = false
  }
}

onMounted(() => {
  loadDrivers()
})
</script>

<style scoped>
.drivers-view {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.panel.hero {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.panel.hero h2 {
  margin-top: 0;
}

.tips {
  list-style: none;
  padding: 0;
  margin: 0;
  color: var(--muted);
  font-size: 0.9rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.tips li {
  padding-left: 1rem;
  position: relative;
}

.tips li::before {
  content: '•';
  position: absolute;
  left: 0;
  color: var(--accent);
}

.panel {
  background: rgba(14, 14, 16, 0.95);
  border-radius: 0.75rem;
  padding: 1rem 1.25rem;
  border: 1px solid var(--border);
  box-shadow: 0 18px 40px rgba(0, 0, 0, 0.4);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.panel-header h2 {
  margin: 0;
  font-size: 1.1rem;
}

.controls {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.875rem;
}

.controls input {
  background: var(--bg-0);
  border-radius: 0.5rem;
  border: 1px solid var(--border);
  color: var(--text);
  padding: 0.3rem 0.5rem;
  width: 8rem;
}

.controls button {
  background: var(--accent);
  border-radius: 9999px;
  border: none;
  padding: 0.4rem 0.9rem;
  color: #fff;
  font-weight: 500;
  cursor: pointer;
}

.controls button:disabled {
  opacity: 0.6;
  cursor: default;
}

.drivers-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 0.75rem;
}

.driver-card {
  background: var(--panel-soft);
  border-radius: 0.75rem;
  padding: 0.6rem 0.75rem;
  border: 1px solid var(--border);
  cursor: pointer;
  transition:
    transform 0.1s ease,
    border-color 0.1s ease,
    box-shadow 0.1s ease;
}

.driver-card:hover {
  transform: translateY(-1px);
  border-color: var(--accent);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.6);
}

.driver-card--selected {
  border-color: var(--accent-2);
}

.driver-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.driver-header img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.driver-number {
  font-weight: 700;
  color: var(--text);
}

.driver-name {
  font-size: 0.9rem;
}

.driver-subtitle {
  display: block;
  font-size: 0.75rem;
  color: var(--muted);
}

.driver-meta {
  margin-top: 0.3rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.team-pill {
  font-size: 0.75rem;
  padding: 0.1rem 0.5rem;
  border-radius: 9999px;
  border: 1px solid transparent;
  color: #0b0b0d;
  background: #6b7280;
}

.country-code {
  font-size: 0.75rem;
  color: var(--muted);
}

.muted {
  color: var(--muted);
  font-size: 0.85rem;
}

.error {
  color: var(--accent-2);
  font-size: 0.85rem;
  margin-bottom: 0.5rem;
}

.driver-layout {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 2fr);
  gap: 1rem;
}

@media (max-width: 900px) {
  .driver-layout {
    grid-template-columns: minmax(0, 1fr);
  }
}

.card {
  background: var(--panel-soft);
  border-radius: 0.75rem;
  padding: 0.75rem 0.9rem;
  border: 1px solid var(--border);
}

.card h3 {
  margin: 0 0 0.4rem;
  font-size: 0.95rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.6rem;
}

.card-header h3 {
  margin: 0;
}

.card-header button {
  background: var(--accent);
  border-radius: 9999px;
  border: none;
  padding: 0.35rem 0.8rem;
  color: #fff;
  font-weight: 500;
  cursor: pointer;
}

.card-header button:disabled {
  opacity: 0.6;
  cursor: default;
}

.driver-profile .profile-header {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  margin-bottom: 0.6rem;
}

.driver-profile img {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  object-fit: cover;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.driver-profile dl {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.4rem;
  margin: 0;
}

.driver-profile dt {
  font-size: 0.7rem;
  color: var(--muted);
}

.driver-profile dd {
  margin: 0;
  font-weight: 500;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem;
  margin-top: 0.4rem;
}

.metric {
  background: var(--bg-0);
  border-radius: 0.5rem;
  padding: 0.35rem 0.5rem;
  border: 1px solid var(--border);
}

.metric-label {
  display: block;
  font-size: 0.7rem;
  color: var(--muted);
}

.metric-value {
  font-weight: 600;
  font-size: 0.95rem;
}

.gear-list {
  list-style: none;
  padding-left: 0;
  margin: 0.4rem 0 0;
  font-size: 0.85rem;
}

.gear-list li {
  display: flex;
  justify-content: space-between;
  padding: 0.2rem 0;
  border-bottom: 1px dashed rgba(255, 255, 255, 0.12);
}

.gear-list li:last-child {
  border-bottom: none;
}

.gear-count {
  color: var(--accent);
}
</style>
