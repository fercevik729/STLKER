import React from "react";
import { Nav, NavLink, NavMenu } 
    from "./NavbarElements";
  
const Navbar = () => {
  return (
    <Nav>
      <NavMenu>
        <NavLink to="/">
          Home
        </NavLink>
        <NavLink to="/about">
          About
        </NavLink>
        <NavLink to="/contact">
          Contact Us
        </NavLink>
        <NavLink to="/sign-up">
          Sign Up
        </NavLink>
        <NavLink to="/log-in">
          Log In 
        </NavLink>
      </NavMenu>
    </Nav>
  )
}
  
export default Navbar;