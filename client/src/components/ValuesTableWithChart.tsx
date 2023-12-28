import {
  Grid,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
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

export const chartOptions = {
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

export default function ValuesTableWithChart (): React.ReactElement {
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
        <Line options={chartOptions} data={data}/>
        <TableContainer component={Paper}>
          <Table aria-label="accounts table">
            <TableHead>
              <TableRow>
                <TableCell><b>Date</b></TableCell>
                <TableCell><b>Value</b></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {
                account.values.map((value, i) => (
                  <TableRow
                    key={i}
                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                  >
                    <TableCell component="th" scope="row">{new Date(value.CreatedAt).toLocaleString()}</TableCell>
                    <TableCell>{moneyFormatter.format(value.value)}</TableCell>
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
