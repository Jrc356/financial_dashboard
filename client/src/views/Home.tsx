import { Grid } from '@mui/material'
import React from 'react'
import NetworthGraph from '../components/NetworthGraph'
import { type Account, GetAllAccounts } from '../lib/api'
import AccountCard from '../components/AccountCard'

export default function Home (): React.ReactElement {
  const [accounts, setAccounts] = React.useState<Account[]>()

  React.useEffect(() => {
    GetAllAccounts()
      .then((accs) => {
        setAccounts(accs)
      })
      .catch(console.error)
  }, [])

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
        <NetworthGraph/>
        {
          accounts?.map((account, i) => (<AccountCard key={i} account={account}></AccountCard>))
        }
      </Grid>
    </Grid>
  )
}
