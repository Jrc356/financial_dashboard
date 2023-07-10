import React from 'react'
import { API, GetAll, type Account, GetValuesForAccount } from '../../lib/api'
import {
  Box,
  Card,
  CardActionArea,
  CardContent,
  Divider,
  Grid,
  styled,
  Typography
} from '@mui/material'
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight'

const formatter = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD'
})

const StyledCardActionArea = styled(CardActionArea)`
    .MuiTouchRipple-child {
        background-color: black;
    }
`

interface Props {
  accountType: API
}

export default function AccountView ({
  accountType
}: Props): React.ReactElement {
  const [accounts, setAccounts] = React.useState([] as Account[])
  const [totalValue, setTotalValue] = React.useState(0)

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

  // TODO: This is so gross
  React.useEffect(() => {
    GetAll(api)
      .then((accs) => {
        setAccounts(accs)
      })
      .catch(console.error)
  }, [])

  React.useEffect(() => {
    if (accounts.length === 0) {
      return
    }
    const valuePromises = []
    for (const account of accounts) {
      if (account.Values.length === 0) {
        valuePromises.push(GetValuesForAccount(api, account))
      }
    }
    if (valuePromises.length > 0) {
      Promise.all(valuePromises)
        .then((accs) => {
          setAccounts(accs)
        })
        .catch(console.error)
    }
  }, [accounts])

  React.useEffect(() => {
    if (accounts.length === 0) {
      return
    }
    let v = 0
    for (const account of accounts) {
      if (account.Values.length > 0) {
        v += account.Values[0].Value
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
      overflow={'hidden'}
    >
      <Grid item>
        <Typography variant="h6" color={'black'} marginTop={6}>
          Total Value:
        </Typography>
        <Typography variant="h2" color={'black'}>{formatter.format(totalValue)}</Typography>
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
        >
        </Divider>
        {accounts.map((account, i) => (
          <Card
            key={i}
            sx={{
              backgroundColor: 'whitesmoke',
              height: 120,
              width: 330,
              margin: 2,
              borderRadius: 3
            }}
          >
            <StyledCardActionArea onClick={() => { console.log(account.Name) }}>
              <CardContent>
                <Box sx={{ display: 'flex', flexDirection: 'row' }}>
                  <Box sx={{ display: 'flex', flexDirection: 'column', flexGrow: 2 }}>
                    <Typography
                      variant="h4"
                      color={'black'}
                      align="left"
                      fontSize={28}
                    >
                      {account.Name.length > 16
                        ? `${account.Name.slice(0, 16)}...`
                        : account.Name}
                    </Typography>
                    <Typography
                      variant="h4"
                      color={'black'}
                      align="left"
                      fontSize={20}
                      paddingTop={3}
                    >
                      {account.Values.length > 0 &&
                        formatter.format(account.Values[0].Value)}
                    </Typography>
                  </Box>
                  <Box sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    textAlign: 'right',
                    textAlignVertical: 'center'
                  }}>
                    {/* display: 'flex', flexDirection: 'column', flexGrow: 3,  */}
                    <KeyboardArrowRightIcon sx={{
                      color: 'black',
                      flex: 1,
                      justifyContent: 'center',
                      alignItems: 'center'
                    }}/>
                  </Box>
                </Box>
              </CardContent>
            </StyledCardActionArea>
          </Card>
        ))}
      </Grid>
    </Grid>
  )
}
