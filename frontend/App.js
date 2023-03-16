import 'regenerator-runtime/runtime';
import React, {useEffect} from 'react';
import PropTypes from 'prop-types';
import {Route, Routes} from 'react-router-dom';
import Home from '../frontend/pages/Home/Home';
import CustomerDashboard from '../frontend/pages/CustomerDashboard/CustomerDashboard';
import NodeDashboard from '../frontend/pages/NodeDashboard/NodeDashboard';
import { connector } from './connector';

const App = () => {
  useEffect(() => {
    connector.restoreConnection();
  }, []);

  return (
    <Routes>
      <Route
        index
        element={<Home/>}
      />
      <Route
        path="customer-dashboard"
        element={<CustomerDashboard/>}
      />
      <Route
        path="node-dashboard"
        element={<NodeDashboard/>}
      />
    </Routes>
  );
};

App.propTypes = {
  currentUser: PropTypes.shape({
    accountId: PropTypes.string.isRequired,
    balance: PropTypes.string.isRequired
  }),
};

export default App;
