import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import React from 'react'
import { green } from '@mui/material/colors'
import AccountView from './views/AccountView'
import { API } from './lib/api'
import { Box } from '@mui/material'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import BottomNav from './components/BottomNav'

const darkTheme = createTheme({
  palette: {
    mode: 'dark'
  }
})

export default function App (): React.ReactElement {
  return (
      <Box style={{
        backgroundColor: green[300],
        textAlign: 'center'
      }}>
        <ThemeProvider theme={darkTheme}>
          <CssBaseline />
          <BrowserRouter>
            <Routes>
              <Route path="/assets" element={<AccountView accountType={API.Asset}/>} />
              <Route path="/liabilities" element={<AccountView accountType={API.Liability}/>} />
            </Routes>
            <BottomNav />
          </BrowserRouter>
        </ThemeProvider>
      </Box>
  )
}
