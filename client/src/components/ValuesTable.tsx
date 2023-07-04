import axios from 'axios'
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Box,
  Tabs,
  Tab,
  Typography
} from '@mui/material'
import React from 'react'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  Legend
} from 'chart.js'
import { Line } from 'react-chartjs-2'

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

interface Account {
  Name: string
  Values: AccountValue[]
}

interface AccountValue {
  Value: number
  Date: string
}

interface Props {
  accountType: string
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
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  )
}

export default function ValuesTable ({ accountType }: Props): JSX.Element {
  const [tabValue, setTabValue] = React.useState(0)
  const [accounts, setAccounts] = React.useState([] as Account[])

  React.useEffect(() => {
    axios.get(`http://localhost:8080/api/${accountType}`)
      .then((response) => {
        const accs = response.data as Account[]
        for (const a of accs) {
          a.Values = []
        }
        setAccounts(accs)
      })
      .catch(console.error)
  }, [])

  React.useEffect(() => {
    if (accounts.length === 0) { return }
    if (accounts[tabValue].Values.length === 0) {
      axios.get(`http://localhost:8080/api/${accountType}/${accounts[tabValue].Name}/value`)
        .then((response) => {
          const a = accounts.slice()
          a[tabValue].Values = response.data
          setAccounts(a)
        })
        .catch(console.error)
    }
  }, [accounts, tabValue])

  if (accounts.length === 0) return <div><p>hi</p></div>

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

  console.log(accounts)
  for (const value of accounts[tabValue].Values) {
    data.labels.push(new Date(value.Date).toLocaleString())
    data.datasets[0].data.unshift(value.Value)
  }

  const handleChange = (event: React.SyntheticEvent, tabValue: number): void => {
    if (accounts[tabValue].Values.length === 0) {
      axios.get(`http://localhost:8080/api/${accountType}/${accounts[tabValue].Name}/value`)
        .then((response) => {
          const a = accounts.slice()
          a[tabValue].Values = response.data
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
            <Table sx={{ minWidth: 650 }} aria-label="accounts table">
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
