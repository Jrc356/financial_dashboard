import React from 'react'
import './App.css'
// import AssetsTable from './components/AssetsTable'
// import LiabilitiesTable from './components/LiabilitiesTable'
import {
  Grid
} from '@mui/material'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import ValuesTable from './components/ValuesTableWithChart'

const darkTheme = createTheme({
  palette: {
    mode: 'dark'
  }
})

export default class App extends React.Component {
  render (): React.ReactNode {
    return (
      <ThemeProvider theme={darkTheme}>
        <CssBaseline />

        <Grid container justifyContent="center" alignItems="center" direction={'column'}>
          <Grid item xs={6} marginTop={3}>
            <ValuesTable accountType="asset"></ValuesTable>
          </Grid>

          {/* <Grid item xs={6} marginTop={3}>
            <AssetsTable></AssetsTable>
          </Grid>

          <Grid item xs={6} marginTop={3}>
            <LiabilitiesTable></LiabilitiesTable>
          </Grid> */}

          {/* <Grid item xs={6} marginTop={3}>
            <ValuesTable accountType="liability"></ValuesTable>
          </Grid> */}
        </Grid>
      </ThemeProvider>
    )
  }
}
