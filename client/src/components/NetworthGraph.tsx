import React from 'react'
import {
  Divider,
  Paper,
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
import { GetNetWorth } from '../lib/api'
import moneyFormatter from '../lib/formatter'

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

export default function NetworthGraph (): React.ReactElement {
  const [networth, setNetworth] = React.useState<Record<string, number>>()

  React.useEffect(() => {
    GetNetWorth()
      .then((nw) => {
        setNetworth(nw)
      })
      .catch(console.error)
  }, [])

  if (networth == null) {
    return <div></div>
  }

  const data = {
    labels: [] as string[],
    datasets: [{
      fill: true,
      label: 'networth',
      data: [] as number[],
      borderColor: 'rgb(53, 162, 235)',
      backgroundColor: 'rgba(53, 162, 235, 0.5)'
    }]
  }

  for (const entry of Object.entries(networth)) {
    data.labels.push(new Date(entry[0]).toLocaleString())
    data.datasets[0].data.push(entry[1])
  }

  return (
    <div>
      <Typography variant="h6" color={'black'} marginTop={6}>
        Networth
      </Typography>
      <Typography variant="h2" color={'black'}>
        {moneyFormatter.format(data.datasets[0].data[data.datasets[0].data.length - 1])}
      </Typography>
      <Paper sx={{ backgroundColor: 'whitesmoke', borderRadius: 3, padding: 2, width: 330, margin: 2 }}>
        <Line options={chartOptions} data={data}/>
      </Paper>
      <Divider
          sx={{
            marginTop: 6,
            marginBottom: 6,
            border: 2,
            borderColor: 'black',
            opacity: 100,
            marginLeft: '5%',
            marginRight: '5%',
            justifyContent: 'center'
          }}
        />
    </div>
  )
}
