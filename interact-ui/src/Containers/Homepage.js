import * as React from 'react';
import {
    AppBar,
    Container,
    Toolbar,
    Box, 
    Button, 
    Typography, 
    Grid, 
    FormControl, 
    InputAdornment, 
    IconButton, 
    OutlinedInput, 
    InputLabel,
    CssBaseline,
    Snackbar
} from '@mui/material'
import { useTheme, ThemeProvider, createTheme } from '@mui/material/styles';
import MuiAlert from '@mui/material/Alert';
import SendIcon from '@mui/icons-material/Send';
import AddIcon from '@mui/icons-material/Add';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
import { createContext, useContext, useState, useMemo } from 'react';
import { useNavigate } from "react-router-dom";
import { UserContext } from '../UserContext';

const pages = ['Home', 'About'];

const Alert = React.forwardRef(function Alert(props, ref) {
    return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

const Homepage = ()=>{

    /***
     * Theme settings
     */
    const [mode, setMode] = useState('light');
    const colorMode = useMemo(
        () => ({
        toggleColorMode: () => {
            setMode((prevMode) => (prevMode === 'light' ? 'dark' : 'light'));
        },
        }),
        [],
    );

    const theme = useMemo(
        () =>
        createTheme({
            palette: {
            mode,
            },
        }),
        [mode],
    )
    
    /**
     * event handlers
     */
     let navigate = useNavigate();

    const [error, setError] = useState("");
    const [showError, setShowError] = useState(false);
    const { setContextDetails } = useContext(UserContext);

    const handleErrorClose = ()=>{
        setShowError(false)
    }
  
    const createEventHandler = () => {
      let createEventURI = process.env.REACT_APP_SERVER_ADDR + "createEvent";
      fetch(createEventURI, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        // TODO: Add a form to collect and send this body
        body: JSON.stringify({
          hostName: "IAmTheHost",
          eventType: "LivePolls",
        }),
      })
        .then((response) => {
          if (!response.ok) {
            console.log("response not ok");
            setError("Join Event failed: API response not ok");
            setShowError(true)
          }
          return response.json();
        })
        .then((data) => {
          if (data.roomId) {
            // TODO: Add hostId also to context
            setContextDetails({ roomId: data.roomId, userId: data.hostId });
            navigate("/hostView", { replace: true });
          } else {
            setError("Create event falied: API didn't return roomId");
            setShowError(true)
          }
          console.log("response json:", data);
        })
        .catch((error) => {
          console.log("error: ", error)
          setError(error.toString())
          setShowError(true)
        });
    };
    const joinEventHandler = (eventID) => {
      let joinEventURI =
        process.env.REACT_APP_SERVER_ADDR + eventID + "/joinEvent";
      fetch(joinEventURI, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      })
        .then((response) => {
          console.log(response);
          if (!response.ok) {
            console.log("response not ok");
            setError(`Could not find the room id ${eventID}`);
            setShowError(true)
          }
          return response.json();
        })
        .then((data) => {
          if (data.clientId) {
            setContextDetails({
              roomId: eventID,
              userId: data.clientId,
            });
            navigate("/AudienceView", { replace: true });
          }
          console.log("response json:", data);
        })
        .catch((error) => {
            console.log(error)
          if (!error.response) {
            setError(error.toString())
          } else {
            setError(error.response.data.message)
          }
          setShowError(true)
        });
    };




    return(
        <ColorModeContext.Provider value={colorMode}>
            <ThemeProvider theme={theme}>
                <CssBaseline />
                <Box sx={{flexGrow:1}}>
                    <HomepageAppBar/>
                    <Container>
                        <Grid container>
                            <Grid item xs={12} md={6}>
                                <InteractCard/>
                                <CreateAndJoinComponent 
                                    create={createEventHandler} 
                                    join={joinEventHandler}
                                />
                            </Grid>
                        </Grid>
                    </Container>
                    <Snackbar open={showError} autoHideDuration={3000} onClose={handleErrorClose}>
                        <Alert onClose={handleErrorClose} severity="error" sx={{ width: '100%' }}>
                            {error}
                        </Alert>
                    </Snackbar>
                </Box>
            </ThemeProvider>
        </ColorModeContext.Provider> 
    )

}



const ColorModeContext = createContext({ toggleColorMode: () => {} });

const HomepageAppBar = ()=>{
    const theme = useTheme();
    const colorMode = useContext(ColorModeContext)


    return(
        <Box sx={{ flexGrow: 1 }}>
            <AppBar position="static">
                <Toolbar>
                    <Container>
                        <Box sx={{ flexGrow: 1, display: { xs: 'none', sm: 'flex' } }}>
                            {pages.map((page) => (
                            <Button
                                key={page}
                                sx={{ my: 2, color: 'white', display: 'block' }}
                            >
                                {page}
                            </Button>
                            ))}
                        </Box>
                    </Container>
                        <IconButton
                            onClick={colorMode.toggleColorMode}    
                        >
                             {theme.palette.mode === 'dark' ? <Brightness7Icon /> : <Brightness4Icon />}
                        </IconButton>
                </Toolbar>
            </AppBar>
        </Box>
    )
}

//TODO: Change height and description
const InteractCard = ()=>{
    return(
        <Box sx={{flexGrow: 1, my: 5}}>
            <Typography variant='h2' sx={{fontWeight: 'bold'}}>
                Interact
            </Typography>
            <Typography>
                Simple and easy way to conduct live quizzes, see the live responses and many more 
            </Typography>
        </Box>
    )
}

const CreateAndJoinComponent = ({create, join})=>{

    const [value, setValue] = useState('')

    const handleChange = (event) => {
        setValue( event.target.value )
    }

    const handleMouseDown = (event) => {
        event.preventDefault();
    }

    return(

        <Grid container spacing={2}>

            <Grid item >
                <Button 
                    variant='contained' 
                    color='secondary'
                    startIcon={<AddIcon/>}
                    sx={{height: '100%'}}
                    onClick={create}   
                >
                    <Typography sx={{textTransform: 'none'}}>
                        Start new event
                    </Typography>
                </Button>
            </Grid>
            
            <Grid item >
                <FormControl  variant="outlined" color='secondary'>
                    <InputLabel>Join with code</InputLabel>
                    <OutlinedInput
                          type={'text'}
                          value={value}
                          onChange={handleChange}
                        endAdornment={
                            <InputAdornment position="end">
                                <IconButton
                                    onClick={()=>{join(value)}}
                                    onMouseDown={handleMouseDown}
                                    edge="end"
                                    >
                                    <SendIcon color='secondary'/>
                                </IconButton>
                            </InputAdornment>
                        }
                        label="Join with code"
                    />
                </FormControl>
            </Grid>
           

        </Grid>
        
    )
}

export default Homepage