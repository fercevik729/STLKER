import React from 'react';
import Navbar from './components/Navbar';
import Home from "./pages/home"
import SignUp from "./pages/signup"
import Contact from "./pages/contact"
import About from "./pages/about"
import { BrowserRouter as _, Routes, Route}
    from 'react-router-dom';

export default function App() {
  return (
    <>
      <Navbar>
      </Navbar>
      <Routes>
        <Route path='/' element={<Home/>}></Route>
        <Route path='/about' element={<About/>}></Route>
        <Route path='/sign-up' element={<SignUp/>}></Route>
        <Route path='/contact' element={<Contact/>}></Route>
      </Routes>
    </>
  )
}

