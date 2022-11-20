import React from "react";
import { Routes, Route } from "react-router-dom";
import Home from "./Home";
import Signup from "./Signup";
import Login from "./Login";
import Verify from "./Verify";
import Portfolio from "./Portfolio";
import Transactions from "./Transactions";
import Health from "./Health";

const Main = () => (
    <Routes>
        <Route exact path="/" element={<Home/>} />
        <Route path="/portfolio/:portfolioId" element={<Portfolio/>} />
        <Route exact path="/signup" element={<Signup/>} />
        <Route exact path="/login" element={<Login/>} />
        <Route exact path="/verify" element={<Verify/>} />
        <Route exact path="/transactions/:portfolioId" element={<Transactions/>} />
        <Route exact path="/health" element={<Health/>} />
    </Routes>
);

export default Main;
