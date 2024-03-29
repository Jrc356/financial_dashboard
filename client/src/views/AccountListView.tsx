import React from 'react'
import { type Account, GetAllAccountsByClass } from '../lib/api'
import {
  Divider,
  Grid,
  Typography
} from '@mui/material'
import AccountCard from '../components/AccountCard'
import moneyFormatter from '../lib/formatter'

function toTitleCase (str: string): string {
  return str.charAt(0).toUpperCase() + str.substring(1).toLowerCase()
}

interface Props {
  accountClass: string
}

export default function AccountListView ({ accountClass }: Props): React.ReactElement {
  const [accounts, setAccounts] = React.useState<Account[]>()
  const [totalValue, setTotalValue] = React.useState(0)

  React.useEffect(() => {
    GetAllAccountsByClass(accountClass)
      .then((accs) => {
        setAccounts(accs)
      })
      .catch(console.error)
  }, [])

  React.useEffect(() => {
    if (accounts == null) {
      return
    }
    // TODO: this isn't right
    let v = 0
    for (const account of accounts) {
      if (account.values.length > 0) {
        console.log(account.values[0].value)
        console.log(typeof account.values[0].value)
        v += Number(account.values[0].value)
      }
    }
    setTotalValue(v)
  }, [accounts])

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
        <Typography variant="h6" color={'black'} marginTop={6}>
          Total {toTitleCase(accountClass)} Value:
        </Typography>
        <Typography variant="h3" color={'black'}>{moneyFormatter.format(totalValue)}</Typography>
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
        {
          accounts?.map((account, i) => (<AccountCard key={i} account={account}></AccountCard>))
        }
      </Grid>
    </Grid>
  )
}
