import api from './client'

export async function fetchDrivers({ sessionKey, meetingKey, driverNumber } = {}) {
  const params = {}
  if (sessionKey) params.session_key = sessionKey
  if (meetingKey) params.meeting_key = meetingKey
  if (driverNumber) params.driver_number = driverNumber

  const res = await api.get('/api/drivers', { params })
  return res.data
}

export async function fetchDriverSpeedSummary({ driverNumber, sessionKey, minSpeed } = {}) {
  const params = {}
  if (driverNumber) params.driver_number = driverNumber
  if (sessionKey) params.session_key = sessionKey
  if (minSpeed) params.min_speed = minSpeed

  const res = await api.get('/api/drivers/speed-summary', { params })
  return res.data
}

export async function fetchCarData({ driverNumber, sessionKey, minSpeed } = {}) {
  const params = {}
  if (driverNumber) params.driver_number = driverNumber
  if (sessionKey) params.session_key = sessionKey
  if (minSpeed) params.min_speed = minSpeed

  const res = await api.get('/api/car-data/raw', { params })
  return res.data
}

export async function fetchGearUsage({ driverNumber, sessionKey, minSpeed } = {}) {
  const params = {}
  if (driverNumber) params.driver_number = driverNumber
  if (sessionKey) params.session_key = sessionKey
  if (minSpeed) params.min_speed = minSpeed

  const res = await api.get('/api/car-data/gear-usage', { params })
  return res.data
}
