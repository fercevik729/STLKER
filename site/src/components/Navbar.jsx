import React from "react";
import { Nav, NavLink, NavMenu } 
    from "./NavbarElements";
  
const Navbar = () => {
  return (
    <Nav>
      <NavMenu>
        <NavLink to="/about">
          About
        </NavLink>
        <NavLink to="/contact">
          Contact Us
        </NavLink>
        <NavLink to="/sign-up">
          Sign Up
        </NavLink>
      </NavMenu>
    </Nav>
  )
}
  
export default Navbar;