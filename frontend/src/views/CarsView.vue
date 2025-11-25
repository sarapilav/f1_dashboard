<template>
  <div class="cars-view">
    <section class="panel">
      <header class="panel-header">
        <div>
          <h2>Telemetry explorer</h2>
          <p class="muted">
            Supply a session and driver number to pull raw telemetry, speed stats and gear usage.
          </p>
        </div>
        <div class="controls">
          <label>
            Session key
            <input v-model="filters.sessionKey" type="number" min="0" />
          </label>
          <label>
            Driver #
            <input v-model="filters.driverNumber" type="number" min="1" />
          </label>
          <label>
            Min speed (km/h)
            <input v-model="filters.minSpeed" type="number" min="0" step="1" />
          </label>
          <button @click="loadTelemetry" :disabled="loading">
            {{ loading ? 'Loading...' : 'Fetch data' }}
          </button>
        </div>
      </header>
      <p v-if="formError" class="error">{{ formError }}</p>
    </section>

    <section class="panel split" v-if="!formError">
      <article>
        <h3>Speed summary</h3>
        <p v-if="summaryError" class="error">{{ summaryError }}</p>
        <div v-else-if="speedSummary" class="stat-grid">
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
        <p v-else-if="!loading" class="muted">Run a query to populate the metrics.</p>
      </article>

      <article>
        <h3>Gear distribution</h3>
        <p v-if="gearError" class="error">{{ gearError }}</p>
        <ul v-else-if="gearUsage.length" class="gear-list">
          <li v-for="item in sortedGearUsage" :key="item.gear">
            <span>Gear {{ item.gear }}</span>
            <span class="gear-count">{{ item.count }}</span>
          </li>
        </ul>
        <p v-else-if="!loading" class="muted">We will render counts per gear after you fetch data.</p>
      </article>
    </section>

    <section class="panel">
      <header class="panel-header">
        <div>
          <h3>Raw car data</h3>
          <p class="muted">
            Showing up to {{ previewData.length }} of {{ telemetry.length }} samples for quick review.
          </p>
        </div>
        <button class="small" @click="copyJSON" :disabled="!telemetry.length">
          Copy JSON
        </button>
      </header>
      <p v-if="telemetryError" class="error">{{ telemetryError }}</p>
      <div v-else-if="previewData.length" class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Date</th>
              <th>Speed</th>
              <th>Gear</th>
              <th>RPM</th>
              <th>Throttle</th>
              <th>Brake</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in previewData" :key="item.date + item.rpm">
              <td>{{ formatDate(item.date) }}</td>
              <td>{{ item.speed.toFixed(1) }} km/h</td>
              <td>{{ item.gear }}</td>
              <td>{{ item.rpm.toFixed(0) }}</td>
              <td>{{ (item.throttle * 100).toFixed(0) }}%</td>
              <td>{{ (item.brake * 100).toFixed(0) }}%</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-else-if="!loading" class="muted">
        No telemetry yet. Enter a session + driver to generate the table.
      </p>
    </section>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { fetchCarData, fetchDriverSpeedSummary, fetchGearUsage } from '../api/openf1'

const filters = reactive({
  sessionKey: '9159',
  driverNumber: '55',
  minSpeed: '',
})

const telemetry = ref([])
const telemetryError = ref('')
const loading = ref(false)

const speedSummary = ref(null)
const summaryError = ref('')

const gearUsage = ref([])
const gearError = ref('')

const formError = ref('')

const previewData = computed(() => telemetry.value.slice(0, 50))
const sortedGearUsage = computed(() => [...gearUsage.value].sort((a, b) => a.gear - b.gear))

function formatDate(date) {
  if (!date) return '—'
  try {
    return new Date(date).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  } catch {
    return date
  }
}

function normalizeNumber(value, label, allowEmpty = false) {
  if (value === '' || value === null || value === undefined) {
    if (allowEmpty) return undefined
    throw new Error(`Missing ${label}.`)
  }
  const parsed = Number(value)
  if (Number.isNaN(parsed)) {
    throw new Error(`${label} must be numeric.`)
  }
  return parsed
}

async function loadTelemetry() {
  telemetryError.value = ''
  summaryError.value = ''
  gearError.value = ''
  formError.value = ''
  loading.value = true
  telemetry.value = []
  speedSummary.value = null
  gearUsage.value = []

  try {
    const sessionKey = normalizeNumber(filters.sessionKey, 'session key')
    const driverNumber = normalizeNumber(filters.driverNumber, 'driver number')
    const minSpeed =
      filters.minSpeed === '' ? undefined : normalizeNumber(filters.minSpeed, 'min speed', true)

    const params = { sessionKey, driverNumber, minSpeed }
    const [raw, summary, gears] = await Promise.all([
      fetchCarData(params),
      fetchDriverSpeedSummary(params),
      fetchGearUsage(params),
    ])

    telemetry.value = raw
    speedSummary.value = summary
    gearUsage.value = gears
  } catch (err) {
    console.error(err)
    if (!speedSummary.value) summaryError.value = 'Failed to load speed summary'
    if (!gearUsage.value.length) gearError.value = 'Failed to load gear usage'
    telemetryError.value = err.message || 'Failed to load car data'
    formError.value = err.message
  } finally {
    loading.value = false
  }
}

async function copyJSON() {
  if (!telemetry.value.length) return
  try {
    await navigator.clipboard.writeText(JSON.stringify(telemetry.value, null, 2))
  } catch (err) {
    console.error(err)
  }
}
</script>

<style scoped>
.cars-view {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
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
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
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

.controls {
  display: flex;
  gap: 0.7rem;
  flex-wrap: wrap;
  align-items: flex-end;
  font-size: 0.85rem;
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

.controls button,
.panel-header button.small {
  background: #0ea5e9;
  border-radius: 9999px;
  border: none;
  padding: 0.45rem 0.95rem;
  color: #0b1120;
  font-weight: 600;
  cursor: pointer;
}

.panel-header button.small {
  align-self: flex-start;
  font-size: 0.85rem;
}

.controls button:disabled,
.panel-header button.small:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.muted {
  color: #94a3b8;
}

.error {
  color: #f87171;
  font-size: 0.9rem;
  margin: 0.2rem 0;
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.7rem;
}

.stat-grid div {
  background: #020617;
  border-radius: 0.75rem;
  border: 1px solid #111827;
  padding: 0.8rem;
}

.stat-grid span {
  color: #94a3b8;
  font-size: 0.8rem;
}

.stat-grid strong {
  font-size: 1.2rem;
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

.table-wrapper {
  overflow-x: auto;
  border: 1px solid #111827;
  border-radius: 0.75rem;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

th,
td {
  padding: 0.5rem;
  text-align: left;
}

thead {
  background: rgba(15, 23, 42, 0.8);
}

tbody tr:nth-child(even) {
  background: rgba(2, 6, 23, 0.3);
}

tbody tr:nth-child(odd) {
  background: rgba(15, 23, 42, 0.4);
}
</style>
