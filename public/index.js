import React from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import Header from "./components/Header";
import Main from "./components/Main";

const root = createRoot(document.getElementById('root'));

root.render(
  <BrowserRouter>
      <Header />
      <Main />
  </BrowserRouter>
);