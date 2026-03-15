import { useState } from 'react'
import './App.css'

const API = 'http://localhost:8080/api/v1'

function App() {

  const [userId, setUserId] = useState('')
  const [action, setAction] = useState('')
  const [metadata, setMetadata] = useState('')

  const [createMessage, setCreateMessage] = useState('')

  async function handleCreate(event) {
    event.preventDefault()

    const response = await fetch(API + '/events', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        user_id: Number(userId),
        action: action,
        metadata: metadata ? JSON.parse(metadata) : null
      })
    })

    const data = await response.json()

    if (response.ok) {
      setCreateMessage('✅ Event created! ID: ' + data.id)
      setUserId('')
      setAction('')
      setMetadata('')
    } else {
      setCreateMessage('❌ Error: ' + data.error)
    }
  }


  const [searchUserId, setSearchUserId] = useState('')
  const [startDate, setStartDate] = useState('')
  const [endDate, setEndDate] = useState('')

  const [events, setEvents] = useState([])

  async function handleSearch(event) {
    event.preventDefault()

    let url = API + '/events?user_id=' + searchUserId
    if (startDate) url += '&start_date=' + new Date(startDate).toISOString()
    if (endDate) url += '&end_date=' + new Date(endDate).toISOString()

    const response = await fetch(url)
    const data = await response.json()

    if (response.ok) {
      setEvents(data || [])
    } else {
      alert('Error: ' + data.error)
    }
  }

  return (
    <div className="container">
      <header className="main-header">
        <h1>Activity Tracker</h1>
        <p>Dashboard for monitoring user events and statistics</p>
      </header>

      <div className="dashboard-grid">
        {/* Форма створення події */}
        <section className="card card-create">
          <div className="card-header">
            <h2>Track New Event</h2>
          </div>
          <form onSubmit={handleCreate} className="event-form">
            <div className="input-group">
              <label>User ID</label>
              <input
                type="number"
                placeholder="Ex: 42"
                value={userId}
                onChange={e => setUserId(e.target.value)}
                required
              />
            </div>
            <div className="input-group">
              <label>Action Type</label>
              <input
                type="text"
                placeholder='Ex: page_view'
                value={action}
                onChange={e => setAction(e.target.value)}
                required
              />
            </div>
            <div className="input-group">
              <label>Metadata (JSON)</label>
              <input
                type="text"
                placeholder='Ex: {"page": "/home"}'
                value={metadata}
                onChange={e => setMetadata(e.target.value)}
              />
            </div>
            <button type="submit" className="button-primary">Send Activity</button>
          </form>

          {createMessage && (
            <div className={`status-banner ${createMessage.startsWith('✅') ? 'success' : 'error'}`}>
              {createMessage}
            </div>
          )}
        </section>

        {/* Форма пошуку подій */}
        <section className="card card-search">
          <div className="card-header">
            <h2>History & Search</h2>
          </div>
          <form onSubmit={handleSearch} className="search-form">
            <div className="input-row">
              <div className="input-group">
                <label>Target User ID</label>
                <input
                  type="number"
                  placeholder="UserID"
                  value={searchUserId}
                  onChange={e => setSearchUserId(e.target.value)}
                  required
                />
              </div>
            </div>
            <div className="input-row">
              <div className="input-group">
                <label>From Date</label>
                <input
                  type="datetime-local"
                  value={startDate}
                  onChange={e => setStartDate(e.target.value)}
                />
              </div>
              <div className="input-group">
                <label>To Date</label>
                <input
                  type="datetime-local"
                  value={endDate}
                  onChange={e => setEndDate(e.target.value)}
                />
              </div>
            </div>
            <button type="submit" className="button-secondary">Fetch History</button>
          </form>

          {events.length > 0 && (
            <div className="table-container">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Action</th>
                    <th>Metadata</th>
                    <th>Timestamp</th>
                  </tr>
                </thead>
                <tbody>
                  {events.map(ev => (
                    <tr key={ev.id}>
                      <td>#{ev.id}</td>
                      <td><span className="action-tag">{ev.action}</span></td>
                      <td className="metadata-cell">{JSON.stringify(ev.metadata)}</td>
                      <td className="date-cell">{new Date(ev.created_at).toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </div>
  )
}

export default App
