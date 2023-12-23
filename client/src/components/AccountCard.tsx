import React from 'react'
import {
  Box,
  Card,
  CardActionArea,
  CardContent,
  Typography
} from '@mui/material'
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight'
import { type Account } from '../lib/api'
import moneyFormatter from '../lib/formatter'
import { Link } from 'react-router-dom'

interface Props {
  account: Account
}

export default function AccountCard ({ account }: Props): React.ReactElement {
  return (
    <Card
      sx={{
        backgroundColor: 'whitesmoke',
        height: 120,
        width: 330,
        margin: 2,
        borderRadius: 3
      }}
    >
      <CardActionArea
        component={Link}
        to={`/accounts/${account.Name}`}
        sx={{
          '&& .MuiTouchRipple-child': {
            backgroundColor: 'black'
          }
        }}>
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
                  moneyFormatter.format(account.Values[0].Value)}
              </Typography>
            </Box>
            <Box sx={{
              display: 'flex',
              flexDirection: 'column',
              textAlign: 'right',
              textAlignVertical: 'center'
            }}>
              <KeyboardArrowRightIcon sx={{
                color: 'black',
                flex: 1,
                justifyContent: 'center',
                alignItems: 'center'
              }}/>
            </Box>
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  )
}
