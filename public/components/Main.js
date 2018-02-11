import React from 'react'
import { Switch, Route } from 'react-router-dom'
import Home from './Home'
import Signup from './Signup'
import Login from './Login'
import Verify from './Verify'
import Portfolio from './Portfolio'

const Main = () => (
  <main>
    <Switch>
      <Route exact path='/' component={Home}/>
      <Route path='/portfolio/:id' component={Portfolio}/>
      <Route exact path='/signup' component={Signup}/>
      <Route exact path='/login' component={Login}/>
      <Route exact path='/verify' component={Verify}/>
    </Switch>
  </main>
)

export default Main