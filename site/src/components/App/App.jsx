import React from 'react';
import Navbar from '../Navbar/Navbar';
import LogIn from '../LogIn/LogIn'
import Home from "../../views/Home"
import SignUp from "../../views/SignUp"
import Contact from "../../views/Contact"
import About from "../../views/About"
import { BrowserRouter as _, Routes, Route}
    from 'react-router-dom';

export default function App() {
  return (
    <>
      <Navbar/>
      <Routes>
        <Route path='/' element={<Home/>}></Route>
        <Route path='/about' element={<About/>}></Route>
        <Route path='/sign-up' element={<SignUp/>}></Route>
        <Route path='/contact' element={<Contact/>}></Route>
        <Route path='/log-in' element={<LogIn/>}></Route>
      </Routes>
    </>
  )
}

