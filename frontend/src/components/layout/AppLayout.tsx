import { Outlet } from 'react-router-dom'

function AppLayout() {
  return (
    <div>
      <div>Sidebar goes here</div>      
      <Outlet />       
    </div>
  )
}

export default AppLayout