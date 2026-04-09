import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import client from '../api/client'
import { useAuthStore } from '../store/authStore'


function Login() {
    const navigate = useNavigate()
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const { setToken } = useAuthStore()
    
    const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
        e.preventDefault()
        try {
             const response = await client.post('/auth/login', {
            Username: username,
            Password: password,
        })
        setToken(response.data.token)
        navigate('/app/dashboard')
        } catch (err) {
            setError('invalid credentials')
        }
    }
     return (
        <div>
        <h1>Login</h1>
        <form onSubmit={handleSubmit}>
            <input
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Username"
            />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
        />
        {error && <p>{error}</p>}
        <button type="submit">Login</button>
      </form>
    </div>
  )
}

export default Login