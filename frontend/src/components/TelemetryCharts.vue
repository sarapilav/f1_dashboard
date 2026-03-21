<template>
  <div class="telemetry-charts">
    <div class="chart-card">
      <header>
        <h3>Speed vs time</h3>
        <span v-if="speedRange" class="range">
          {{ speedRange.min.toFixed(0) }}–{{ speedRange.max.toFixed(0) }} km/h
        </span>
      </header>
      <div v-if="speedPath" class="chart">
        <svg viewBox="0 0 1000 220" preserveAspectRatio="none">
          <defs>
            <linearGradient id="speedFill" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="var(--accent)" stop-opacity="0.45" />
              <stop offset="100%" stop-color="var(--accent)" stop-opacity="0" />
            </linearGradient>
          </defs>
          <path :d="speedAreaPath" fill="url(#speedFill)" />
          <path :d="speedPath" fill="none" stroke="var(--accent)" stroke-width="2.2" />
        </svg>
        <div class="axis axis-y">Speed (km/h)</div>
        <div class="axis axis-x">Time</div>
      </div>
      <p v-else class="muted">Fetch telemetry to render the speed chart.</p>
    </div>

    <div class="chart-card">
      <header>
        <h3>Gear over time</h3>
        <span v-if="gearRange" class="range">G{{ gearRange.min }}–G{{ gearRange.max }}</span>
      </header>
      <div v-if="gearPath" class="chart">
        <svg viewBox="0 0 1000 220" preserveAspectRatio="none">
          <path :d="gearPath" fill="none" stroke="var(--accent-2)" stroke-width="2.2" />
        </svg>
        <div class="axis axis-y">Gear</div>
        <div class="axis axis-x">Time</div>
      </div>
      <p v-else class="muted">Fetch telemetry to render the gear chart.</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  telemetry: {
    type: Array,
    default: () => [],
  },
})

const MAX_POINTS = 350
const WIDTH = 1000
const HEIGHT = 220
const PAD_X = 24
const PAD_Y = 20

const normalized = computed(() => {
  if (!props.telemetry || props.telemetry.length < 2) return []
  const sorted = [...props.telemetry]
    .map((d) => ({
      t: Date.parse(d.date),
      speed: Number(d.speed),
      gear: Number(d.n_gear ?? d.gear),
    }))
    .filter((d) => Number.isFinite(d.t) && Number.isFinite(d.speed) && Number.isFinite(d.gear))
    .sort((a, b) => a.t - b.t)

  if (sorted.length <= MAX_POINTS) return sorted
  const stride = Math.ceil(sorted.length / MAX_POINTS)
  return sorted.filter((_, idx) => idx % stride === 0)
})

const speedRange = computed(() => {
  if (!normalized.value.length) return null
  let min = Infinity
  let max = -Infinity
  normalized.value.forEach((d) => {
    if (d.speed < min) min = d.speed
    if (d.speed > max) max = d.speed
  })
  return { min, max }
})

const gearRange = computed(() => {
  if (!normalized.value.length) return null
  let min = Infinity
  let max = -Infinity
  normalized.value.forEach((d) => {
    if (d.gear < min) min = d.gear
    if (d.gear > max) max = d.gear
  })
  return { min, max }
})

function scaleX(t, minT, maxT) {
  if (maxT === minT) return PAD_X
  return PAD_X + ((t - minT) / (maxT - minT)) * (WIDTH - PAD_X * 2)
}

function scaleY(v, minV, maxV) {
  if (maxV === minV) return HEIGHT / 2
  const ratio = (v - minV) / (maxV - minV)
  return HEIGHT - PAD_Y - ratio * (HEIGHT - PAD_Y * 2)
}

const speedPath = computed(() => {
  if (!normalized.value.length || !speedRange.value) return ''
  const minT = normalized.value[0].t
  const maxT = normalized.value[normalized.value.length - 1].t
  const { min, max } = speedRange.value

  const parts = normalized.value.map((d, idx) => {
    const x = scaleX(d.t, minT, maxT)
    const y = scaleY(d.speed, min, max)
    return `${idx === 0 ? 'M' : 'L'} ${x.toFixed(2)} ${y.toFixed(2)}`
  })
  return parts.join(' ')
})

const speedAreaPath = computed(() => {
  if (!speedPath.value) return ''
  const minT = normalized.value[0].t
  const maxT = normalized.value[normalized.value.length - 1].t
  const baseY = HEIGHT - PAD_Y
  const firstX = scaleX(minT, minT, maxT)
  const lastX = scaleX(maxT, minT, maxT)
  return `${speedPath.value} L ${lastX.toFixed(2)} ${baseY} L ${firstX.toFixed(2)} ${baseY} Z`
})

const gearPath = computed(() => {
  if (!normalized.value.length || !gearRange.value) return ''
  const minT = normalized.value[0].t
  const maxT = normalized.value[normalized.value.length - 1].t
  const { min, max } = gearRange.value

  const points = normalized.value.map((d) => ({
    x: scaleX(d.t, minT, maxT),
    y: scaleY(d.gear, min, max),
  }))

  let d = `M ${points[0].x.toFixed(2)} ${points[0].y.toFixed(2)}`
  for (let i = 1; i < points.length; i += 1) {
    d += ` H ${points[i].x.toFixed(2)} V ${points[i].y.toFixed(2)}`
  }
  return d
})
</script>

<style scoped>
.telemetry-charts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
}

.chart-card {
  background: var(--panel-soft);
  border: 1px solid var(--border);
  border-radius: 0.75rem;
  padding: 0.9rem;
}

.chart-card header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 0.75rem;
  margin-bottom: 0.6rem;
}

.chart-card h3 {
  margin: 0;
  font-size: 0.95rem;
}

.range {
  color: var(--muted);
  font-size: 0.75rem;
}

.chart {
  position: relative;
  height: 200px;
  background: var(--bg-0);
  border-radius: 0.6rem;
  border: 1px solid var(--border);
  padding: 0.4rem;
}

.chart svg {
  width: 100%;
  height: 100%;
}

.axis {
  position: absolute;
  color: var(--muted);
  font-size: 0.7rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.axis-y {
  top: 50%;
  left: 0.35rem;
  transform: translateY(-50%) rotate(-90deg);
  transform-origin: left center;
}

.axis-x {
  bottom: 0.35rem;
  right: 0.75rem;
}

.muted {
  color: var(--muted);
  font-size: 0.85rem;
}
</style>
