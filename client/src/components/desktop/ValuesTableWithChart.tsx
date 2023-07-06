import {
  Box,
  Paper,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tabs
} from '@mui/material'
import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  Title,
  Tooltip
} from 'chart.js'
import React from 'react'
import { Line } from 'react-chartjs-2'
import { API, GetAll, GetValuesForAccount, type Account } from '../../lib/api'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  Legend
)

export const chartOptions = {
  responsive: true,
  plugins: {
    legend: {
      position: 'top' as const
    },
    title: {
      display: true,
      text: 'Value over Time'
    }
  }
}

const formatter = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD'
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
          {children}
        </Box>
      )}
    </div>
  )
}

interface Props {
  accountType: string
}

export default function ValuesTableWithChart ({ accountType }: Props): React.ReactElement {
  const [tabValue, setTabValue] = React.useState(0)
  const [accounts, setAccounts] = React.useState([] as Account[])

  let api = API.Asset
  switch (accountType) {
    case API.Asset:
      api = API.Asset
      break
    case API.Liability:
      api = API.Liability
      break
    default:
      console.error('unknown account type')
      break
  }

  React.useEffect(() => {
    GetAll(api)
      .then((accs) => {
        setAccounts(accs)
      })
      .catch(console.error)
  }, [])

  React.useEffect(() => {
    if (accounts.length === 0) { return }
    if (accounts[tabValue].Values.length === 0) {
      GetValuesForAccount(api, accounts[tabValue])
        .then((acc) => {
          const a = accounts.slice()
          a[tabValue] = acc
          setAccounts(a)
        })
        .catch(console.error)
    }
  }, [accounts, tabValue])

  if (accounts.length === 0) return <div></div>

  const data = {
    labels: [] as string[],
    datasets: [{
      fill: true,
      label: accounts[tabValue].Name,
      data: [] as number[],
      borderColor: 'rgb(53, 162, 235)',
      backgroundColor: 'rgba(53, 162, 235, 0.5)'
    }]
  }

  for (const value of accounts[tabValue].Values) {
    // TODO: drop time from date
    // TODO: depends on having unique days
    data.labels.push(new Date(value.Date).toLocaleString())
    data.datasets[0].data.unshift(value.Value)
  }

  const handleChange = (event: React.SyntheticEvent, tabValue: number): void => {
    if (accounts[tabValue].Values.length === 0) {
      GetValuesForAccount(api, accounts[tabValue])
        .then((acc) => {
          const a = accounts.slice()
          a[tabValue] = acc
          setAccounts(a)
        })
        .catch(console.error)
    }
    setTabValue(tabValue)
  }

  return (
    <div>
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={handleChange} aria-label="accounts">
        {accounts.map((account) => (
          <Tab key={account.Name} label={account.Name} />
        ))}
        </Tabs>
      </Box>

      {accounts.map((account, i) => (
        <TabPanel key={account.Name} value={tabValue} index={i}>
          <Line options={chartOptions} data={data} />
          <TableContainer component={Paper}>
            <Table aria-label="accounts table">
              <TableHead>
                <TableRow>
                  <TableCell><b>Date</b></TableCell>
                  <TableCell><b>Value</b></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {accounts[i].Values.map((value) => (
                  <TableRow
                    key={Date.parse(value.Date).toString()}
                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                  >
                    <TableCell component="th" scope="row">{new Date(value.Date).toLocaleString()}</TableCell>
                    <TableCell>{formatter.format(value.Value)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </TabPanel>
      ))}
    </div>
  )
}
