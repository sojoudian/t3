import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Card, Form, Button, Alert, Spinner } from 'react-bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';

function App() {
  const [currentTimes, setCurrentTimes] = useState({
    torontoTime: '',
    tehranTime: '',
    torontoTimeStr: '',
    tehranTimeStr: ''
  });
  
  const [conversionData, setConversionData] = useState({
    sourceCity: 'Toronto',
    hour: 12,
    minute: 0,
    result: null
  });
  
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  
  // Helper function to get the base API URL
  const getApiBaseUrl = () => {
    // Check if we're running in development or production
    if (process.env.NODE_ENV === 'development') {
      return 'http://localhost:8080'; // Development backend URL
    }
    return ''; // In production, use relative URLs
  };
  
  // Fetch current times on component mount
  useEffect(() => {
    fetchCurrentTimes();
    
    // Refresh times every minute
    const intervalId = setInterval(fetchCurrentTimes, 60000);
    
    // Clean up the interval on component unmount
    return () => clearInterval(intervalId);
  }, []);
  
  const fetchCurrentTimes = async () => {
    try {
      setLoading(true);
      console.log('Fetching current times...');
      
      const apiUrl = `${getApiBaseUrl()}/api/current-time`;
      console.log('API URL:', apiUrl);
      
      const response = await fetch(apiUrl);
      
      if (!response.ok) {
        throw new Error(`Failed to fetch current times: ${response.status} ${response.statusText}`);
      }
      
      const data = await response.json();
      console.log('Received time data:', data);
      
      setCurrentTimes({
        torontoTime: data.toronto_time || '',
        tehranTime: data.tehran_time || '',
        torontoTimeStr: data.toronto_time_str || 'Time not available',
        tehranTimeStr: data.tehran_time_str || 'Time not available'
      });
      
      setError(null);
    } catch (err) {
      console.error('Error fetching current times:', err);
      setError(`Failed to fetch current times: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };
  
  const handleCityChange = (e) => {
    setConversionData({
      ...conversionData,
      sourceCity: e.target.value
    });
  };
  
  const handleHourChange = (e) => {
    setConversionData({
      ...conversionData,
      hour: parseInt(e.target.value, 10)
    });
  };
  
  const handleMinuteChange = (e) => {
    setConversionData({
      ...conversionData,
      minute: parseInt(e.target.value, 10)
    });
  };
  
  const handleConvert = async () => {
    try {
      setError(null);
      console.log('Converting time...');
      
      const apiUrl = `${getApiBaseUrl()}/api/convert-time`;
      console.log('API URL:', apiUrl);
      console.log('Request data:', {
        city: conversionData.sourceCity,
        hour: conversionData.hour,
        minute: conversionData.minute
      });
      
      const response = await fetch(apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          city: conversionData.sourceCity,
          hour: conversionData.hour,
          minute: conversionData.minute
        })
      });
      
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Failed to convert time: ${response.status} ${response.statusText} - ${errorText}`);
      }
      
      const data = await response.json();
      console.log('Received conversion data:', data);
      
      setConversionData({
        ...conversionData,
        result: data
      });
    } catch (err) {
      console.error('Error converting time:', err);
      setError(`Failed to convert time: ${err.message}`);
    }
  };
  
  // Generate hours for dropdown (0-23)
  const hours = Array.from({ length: 24 }, (_, i) => i);
  
  // Generate minutes for dropdown (0, 15, 30, 45)
  const minutes = [0, 15, 30, 45];
  
  return (
    <Container className="py-5">
      <h1 className="text-center mb-5">Time Converter: Toronto & Tehran</h1>
      
      {error && <Alert variant="danger">{error}</Alert>}
      
      <Row className="mb-5">
        <Col md={6} className="mb-3">
          <Card className="h-100">
            <Card.Header className="bg-primary text-white">
              <h2 className="h5 mb-0">Toronto Time</h2>
            </Card.Header>
            <Card.Body className="d-flex align-items-center justify-content-center">
              {loading ? (
                <Spinner animation="border" variant="primary" />
              ) : (
                <h3 className="text-center">{currentTimes.torontoTimeStr}</h3>
              )}
            </Card.Body>
          </Card>
        </Col>
        <Col md={6} className="mb-3">
          <Card className="h-100">
            <Card.Header className="bg-success text-white">
              <h2 className="h5 mb-0">Tehran Time</h2>
            </Card.Header>
            <Card.Body className="d-flex align-items-center justify-content-center">
              {loading ? (
                <Spinner animation="border" variant="success" />
              ) : (
                <h3 className="text-center">{currentTimes.tehranTimeStr}</h3>
              )}
            </Card.Body>
          </Card>
        </Col>
      </Row>
      
      <Card className="mb-5">
        <Card.Header className="bg-info text-white">
          <h2 className="h5 mb-0">Convert Time</h2>
        </Card.Header>
        <Card.Body>
          <Form>
            <Row className="align-items-end">
              <Col md={3} className="mb-3">
                <Form.Group>
                  <Form.Label>Select City</Form.Label>
                  <Form.Select 
                    value={conversionData.sourceCity} 
                    onChange={handleCityChange}
                  >
                    <option value="Toronto">Toronto</option>
                    <option value="Tehran">Tehran</option>
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={3} className="mb-3">
                <Form.Group>
                  <Form.Label>Hour</Form.Label>
                  <Form.Select 
                    value={conversionData.hour} 
                    onChange={handleHourChange}
                  >
                    {hours.map(hour => (
                      <option key={hour} value={hour}>
                        {hour.toString().padStart(2, '0')}
                      </option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={3} className="mb-3">
                <Form.Group>
                  <Form.Label>Minute</Form.Label>
                  <Form.Select 
                    value={conversionData.minute} 
                    onChange={handleMinuteChange}
                  >
                    {minutes.map(minute => (
                      <option key={minute} value={minute}>
                        {minute.toString().padStart(2, '0')}
                      </option>
                    ))}
                  </Form.Select>
                </Form.Group>
              </Col>
              <Col md={3} className="mb-3">
                <Button variant="primary" onClick={handleConvert} className="w-100">
                  Convert
                </Button>
              </Col>
            </Row>
          </Form>
          
          {conversionData.result && (
            <Alert variant="success" className="mt-3">
              <h5 className="mb-2">Conversion Result:</h5>
              <p className="mb-0">
                {conversionData.result.source_time} in {conversionData.result.source_city} is{' '}
                <strong>{conversionData.result.target_time}</strong> in {conversionData.result.target_city}
              </p>
            </Alert>
          )}
        </Card.Body>
      </Card>
      
      <Card>
        <Card.Header className="bg-secondary text-white">
          <h2 className="h5 mb-0">About Time Zones</h2>
        </Card.Header>
        <Card.Body>
          <Row>
            <Col md={6}>
              <h5>Toronto (Eastern Time)</h5>
              <p>
                Toronto follows Eastern Time (ET), which can be either Eastern Standard Time (EST, UTC-5)
                or Eastern Daylight Time (EDT, UTC-4) depending on daylight saving time.
              </p>
            </Col>
            <Col md={6}>
              <h5>Tehran (Iran Standard Time)</h5>
              <p>
                Tehran follows Iran Standard Time (IRST, UTC+3:30) or Iran Daylight Time (IRDT, UTC+4:30)
                depending on daylight saving time.
              </p>
            </Col>
          </Row>
        </Card.Body>
      </Card>
    </Container>
  );
}

export default App;
