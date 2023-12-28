import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import React from 'react'
import ReactDOM from 'react-dom/client'
import { blue } from '@mui/material/colors'
import AccountList from './views/AccountList'
import { Box } from '@mui/material'
import {
  createBrowserRouter,
  RouterProvider
} from 'react-router-dom'
import BottomNav from './components/BottomNav'
import Home from './views/Home'
import Account from './views/Account'

function Base (children?: React.ReactElement): React.ReactElement {
  const darkTheme = createTheme({
    palette: {
      mode: 'dark'
    }
  })
  return (
  <Box style={{
    backgroundColor: blue[100],
    textAlign: 'center',
    display: 'flex',
    flexDirection: 'column',
    overflowY: 'auto',
    paddingBottom: '50px',
    paddingTop: '25px'
  }}>
      <ThemeProvider theme={darkTheme}>
        <CssBaseline />
        {children}
        <BottomNav />
      </ThemeProvider>
    </Box>
  )
}

const router = createBrowserRouter([
  {
    path: '/',
    element: Base(<Home/>)
  },
  {
    path: '/assets',
    element: Base(<AccountList accountClass={'asset'}/>)
  },
  {
    path: '/liabilities',
    element: Base(<AccountList accountClass={'liability'}/>)
  },
  {
    path: '/accounts',
    element: Base(<Account/>)
  }
])

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <RouterProvider router={router} />
)
