import React, { type ReactNode } from 'react'
import './App.css'
import AssetsTable from './components/AssetsTable'
import LiabilitiesTable from './components/LiabilitiesTable'
import {
  Grid,
  Box,
  Tabs,
  Tab,
  Typography
} from '@mui/material'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import axios from 'axios'

const darkTheme = createTheme({
  palette: {
    mode: 'dark'
  }
})

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

function TabPanel (props: TabPanelProps): React.ReactElement {
  const { children, value, index, ...other } = props

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  )
}

interface Asset {
  Name: string
  Type: string
  TaxBucket: string
  Values: AssetValue[]
}

interface AssetValue {
  Value: number
  Date: string
}

const formatter = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD'
})

export default class App extends React.Component {
  state = {
    tabValue: 0,
    assets: [] as Asset[]
  }

  componentDidMount (): void {
    axios.get('http://localhost:8080/api/asset')
      .then(res => {
        const assets: Asset = res.data
        this.setState({ assets })
      }).catch(console.error)
  }

  render (): ReactNode {
    const handleChange = (event: React.SyntheticEvent, tabValue: number): void => {
      axios.get(`http://localhost:8080/api/asset/${this.state.assets[tabValue].Name}/value`)
        .then((res) => {
          const { assets } = this.state
          assets[tabValue].Values = res.data as AssetValue[]
          console.log(assets)
          this.setState({ assets, tabValue })
        }).catch(console.error)
    }

    return (
      <ThemeProvider theme={darkTheme}>
        <CssBaseline />

        <Grid container spacing={2} justifyContent="center" alignItems="center">
          <Grid item xs={8}>
            <AssetsTable></AssetsTable>
          </Grid>

          <Grid item xs={8}>
          <LiabilitiesTable></LiabilitiesTable>
          </Grid>

          <Grid item xs={8}>
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
              <Tabs value={this.state.tabValue} onChange={handleChange} aria-label="basic tabs example">
              {this.state.assets.map((asset) => (
                <Tab key={asset.Name} label={asset.Name} />
              ))}
              </Tabs>
            </Box>

            {this.state.assets.map((asset, i) => (
              <TabPanel key={asset.Name} value={this.state.tabValue} index={i}>
                <ul>
                {asset.Values?.map((value, i) => {
                  return <li key={i}>{new Date(value.Date).toLocaleString()}: {formatter.format(value.Value)}</li>
                })}
                </ul>
              </TabPanel>
            ))}
          </Grid>
        </Grid>
      </ThemeProvider>
    )
  }
}
