import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import React from 'react'
import { blue } from '@mui/material/colors'
import AccountView from './components/mobile/AccountView'
import { API } from './lib/api'
import { Box } from '@mui/material'

const darkTheme = createTheme({
  palette: {
    mode: 'dark'
  }
})

export default class App extends React.Component {
  render (): React.ReactNode {
    return (
      <Box style={{
        backgroundColor: blue[800],
        textAlign: 'center'
      }}>
        <ThemeProvider theme={darkTheme}>
          <CssBaseline />
          <AccountView accountType={API.Asset}></AccountView>
        </ThemeProvider>
      </Box>
    )
  }
}
