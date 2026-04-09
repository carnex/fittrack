import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import PublicDashboard from './pages/PublicDashboard'
import AppLayout from './components/layout/AppLayout'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<PublicDashboard />} />
        <Route path="/login" element={<Login />} />
        <Route path="/app" element={<AppLayout />}>
          <Route path="dashboard" element={<Dashboard />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
