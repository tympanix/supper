import React, { Component } from 'react';
import { render } from 'react-dom';

import '../sass/styles.scss';

import { BrowserRouter, Route, Switch, Link } from 'react-router-dom'

import App from './comp/App'

render(<App/>, document.getElementById('root'));
