import React, { Component } from 'react';
import { render } from 'react-dom';

import '../sass/styles.scss';
import 'miqu/styles.scss';
import 'flag-icon-css/sass/flag-icon.scss';

import './websocket.js'

import App from './comp/App'

render(<App/>, document.getElementById('root'));
