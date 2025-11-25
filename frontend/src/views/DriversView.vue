<template>
  <div class="drivers-view">
    <section class="panel hero">
      <div>
        <h2>Driver intelligence</h2>
        <p class="muted">
          Search for a session, list all drivers, then focus on a specific competitor to reveal speed
          summaries, gear usage and raw telemetry samples.
        </p>
      </div>
      <ul class="tips">
        <li>Try session 9159 for a 2023 race.</li>
        <li>You can chain queries quickly from keyboard (Enter must focus button).</li>
        <li>The telemetry preview shows the first five points from /api/car-data/raw.</li>
      </ul>
    </section>

    <section class="panel">
      <div class="panel-header">
        <h2>Drivers</h2>
        <div class="controls">
          <label>
            Session key:
            <input v-model="sessionKeyInput" type="number" placeholder="e.g. 9159" />
          </label>
          <button @click="loadDrivers" :disabled="loadingDrivers">
            {{ loadingDrivers ? 'Loading...' : 'Load drivers' }}
          </button>
        </div>
      </div>

      <div v-if="driversError" class="error">
        {{ driversError }}
      </div>

      <div class="drivers-list" v-if="drivers.length">
        <div
          v-for="driver in drivers"
          :key="driver.driver_number"
          :class="[
            'driver-card',
            driver.driver_number === selectedDriverNumber ? 'driver-card--selected' : '',
          ]"
          @click="selectDriver(driver)"
        >
          <div class="driver-header">
            <span class="driver-number">#{{ driver.driver_number }}</span>
            <span class="driver-name">{{ driver.full_name }}</span>
          </div>
          <div class="driver-meta">
            <span class="team-pill" :style="teamColorStyle(driver.team_colour)">
              {{ driver.team_name || 'Unknown team' }}
            </span>
            <span class="country-code">{{ driver.country_code }}</span>
          </div>
        </div>
      </div>

      <div v-else-if="!loadingDrivers">
        <p class="muted">No drivers loaded yet. Enter a session key and click "Load drivers".</p>
      </div>
    </section>

    <section class="panel" v-if="selectedDriverNumber">
      <div class="panel-header">
        <h2>Driver details</h2>
        <p class="muted">
          Session {{ currentSessionKey }} · Driver #{{ selectedDriverNumber }}
          {{ selectedDriverName }}
        </p>
      </div>

      <div v-if="driverSummaryError" class="error">
        {{ driverSummaryError }}
      </div>

      <div class="driver-layout">
        <div class="summary-cards">
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

        <!-- Placeholder for charts later -->
        <div class="card">
          <h3>Charts placeholder</h3>
          <p class="muted">
            Here you can plug in speed vs time, throttle vs time, etc. Data comes from
            /api/car-data/raw.
          </p>
          <button @click="loadCarData" :disabled="loadingCarData">
            {{ loadingCarData ? 'Loading telemetry...' : 'Fetch raw telemetry' }}
          </button>
          <pre v-if="carData.length" class="telemetry-preview">
            {{ carData.slice(0, 5) }}
          </pre>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { fetchDrivers, fetchDriverSpeedSummary, fetchGearUsage, fetchCarData } from '../api/openf1'

const sessionKeyInput = ref('')
const currentSessionKey = ref(null)
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
    if (!sessionKeyInput.value) {
      throw new Error('Please enter a session_key (e.g. 9159).')
    }

    currentSessionKey.value = Number(sessionKeyInput.value)
    const data = await fetchDrivers({ sessionKey: currentSessionKey.value })

    drivers.value = data
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
      sessionKey: currentSessionKey.value,
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
      sessionKey: currentSessionKey.value,
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
      sessionKey: currentSessionKey.value,
    })
    carData.value = data
  } catch (err) {
    console.error(err)
  } finally {
    loadingCarData.value = false
  }
}
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
  color: #94a3b8;
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
  color: #0ea5e9;
}

.panel {
  background: rgba(15, 23, 42, 0.9);
  border-radius: 0.75rem;
  padding: 1rem 1.25rem;
  border: 1px solid #1f2937;
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
  background: #020617;
  border-radius: 0.5rem;
  border: 1px solid #374151;
  color: #e5e7eb;
  padding: 0.3rem 0.5rem;
  width: 8rem;
}

.controls button {
  background: #0ea5e9;
  border-radius: 9999px;
  border: none;
  padding: 0.4rem 0.9rem;
  color: #0b1120;
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
  background: #020617;
  border-radius: 0.75rem;
  padding: 0.6rem 0.75rem;
  border: 1px solid #1f2937;
  cursor: pointer;
  transition:
    transform 0.1s ease,
    border-color 0.1s ease,
    box-shadow 0.1s ease;
}

.driver-card:hover {
  transform: translateY(-1px);
  border-color: #0ea5e9;
  box-shadow: 0 10px 25px rgba(15, 23, 42, 0.7);
}

.driver-card--selected {
  border-color: #22c55e;
}

.driver-header {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
}

.driver-number {
  font-weight: 700;
  color: #e5e7eb;
}

.driver-name {
  font-size: 0.9rem;
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
  color: #020617;
  background: #6b7280;
}

.country-code {
  font-size: 0.75rem;
  color: #9ca3af;
}

.muted {
  color: #9ca3af;
  font-size: 0.85rem;
}

.error {
  color: #f87171;
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
  background: #020617;
  border-radius: 0.75rem;
  padding: 0.75rem 0.9rem;
  border: 1px solid #1f2937;
}

.card h3 {
  margin: 0 0 0.4rem;
  font-size: 0.95rem;
}

.card-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem;
  margin-top: 0.4rem;
}

.metric {
  background: #020617;
  border-radius: 0.5rem;
  padding: 0.35rem 0.5rem;
  border: 1px solid #111827;
}

.metric-label {
  display: block;
  font-size: 0.7rem;
  color: #9ca3af;
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
  border-bottom: 1px dashed #1f2937;
}

.gear-list li:last-child {
  border-bottom: none;
}

.gear-count {
  color: #a5b4fc;
}

.telemetry-preview {
  background: #020617;
  border-radius: 0.5rem;
  padding: 0.5rem;
  border: 1px solid #111827;
  font-size: 0.7rem;
  max-height: 200px;
  overflow: auto;
  margin-top: 0.5rem;
}
</style>
