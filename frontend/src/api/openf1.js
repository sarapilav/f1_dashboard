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

export async function fetchMeetings({ year, meetingKey, meetingName, location, countryName } = {}) {
  const params = {}
  if (year) params.year = year
  if (meetingKey) params.meeting_key = meetingKey
  if (meetingName) params.meeting_name = meetingName
  if (location) params.location = location
  if (countryName) params.country_name = countryName

  const res = await api.get('/api/meetings', { params })
  return res.data
}

export async function fetchSessions({ meetingKey, sessionKey, sessionType, sessionName } = {}) {
  const params = {}
  if (meetingKey) params.meeting_key = meetingKey
  if (sessionKey) params.session_key = sessionKey
  if (sessionType) params.session_type = sessionType
  if (sessionName) params.session_name = sessionName

  const res = await api.get('/api/sessions', { params })
  return res.data
}

export async function fetchSessionResults({ sessionKey, maxPosition } = {}) {
  const params = {}
  if (sessionKey) params.session_key = sessionKey
  if (maxPosition) params.max_position = maxPosition

  const res = await api.get('/api/session-results', { params })
  return res.data
}
