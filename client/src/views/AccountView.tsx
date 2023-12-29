import React from 'react'
import {
  Grid,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography
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
import { Line } from 'react-chartjs-2'
import { type Account, GetAccountByName } from '../lib/api'
import moneyFormatter from '../lib/formatter'
import { useSearchParams } from 'react-router-dom'

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

const chartOptions = {
  responsive: true,
  plugins: {
    legend: {
      display: false
    },
    title: {
      display: false
    }
  },
  scales: {
    x: {
      display: false
    }
  }
}

export default function AccountView (): React.ReactElement {
  const [account, setAccount] = React.useState<Account>()
  const [searchParams] = useSearchParams()

  const accountName = searchParams.get('name') ?? ''

  React.useEffect(() => {
    GetAccountByName(accountName)
      .then((acc) => {
        setAccount(acc)
      })
      .catch(console.error)
  }, [])

  if (account == null) return <div></div>

  const data = {
    labels: [] as string[],
    datasets: [{
      fill: true,
      label: account.name,
      data: [] as number[],
      borderColor: 'rgb(53, 162, 235)',
      backgroundColor: 'rgba(53, 162, 235, 0.5)'
    }]
  }

  for (const value of account.values) {
    data.labels.unshift(new Date(value.CreatedAt).toLocaleString())
    data.datasets[0].data.unshift(value.value)
  }

  return (
    <Grid
      container
      justifyContent="center"
      alignItems="center"
      direction={'column'}
      overflow={'auto'}
      flex={1}
    >
      <Grid item>
        <Typography variant="h6" color={'black'} marginBottom={1}>
          { account.name }
        </Typography>
        <Paper sx={{ backgroundColor: 'whitesmoke', borderRadius: 3, padding: 2 }}>
          <Line options={chartOptions} data={data}/>
        </Paper>
        <TableContainer
          component={Paper}
          sx={{
            backgroundColor: 'whitesmoke',
            marginTop: 1,
            borderRadius: 3,
            paddingLeft: 2,
            paddingRight: 2
          }}
        >
          <Table aria-label="accounts table">
            <TableHead>
              <TableRow>
                <TableCell sx={{ color: 'black' }}><b>Date</b></TableCell>
                <TableCell sx={{ color: 'black' }}><b>Value</b></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {
                account.values.map((value, i) => (
                  <TableRow
                    key={i}
                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                  >
                    <TableCell component="th" scope="row" sx={{ color: 'black' }}>{new Date(value.CreatedAt).toLocaleString()}</TableCell>
                    <TableCell sx={{ color: 'black' }}>{moneyFormatter.format(value.value)}</TableCell>
                  </TableRow>
                ))
              }
            </TableBody>
          </Table>
        </TableContainer>
      </Grid>
    </Grid>
  )
}
