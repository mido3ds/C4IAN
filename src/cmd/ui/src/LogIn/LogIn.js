import React from 'react';
import {
    withStyles,
    Arwes,
    Content,
    Words
} from 'arwes';
import uImage from '../images/unit.png';
import withTemplate from '../withTemplate';
import TextField from './TextField/TextField';
import './LogIn.css';

const styles = theme => ({
    root: {
        background: 'radial-gradient(circle, rgba(14,63,87,0.9164040616246498) 0%, rgba(0,0,0,0.9472163865546218) 81%)',
    },
});

class LogIn extends React.Component {
    render() {
        const { classes } = this.props;

        return (

            <div>
                {!this.props.port ? <> </> :
                    <Arwes>
                        <Content className={`logIn-root ${classes.root}`}>
                            <Words animate className="back-text"> Back </Words>
                            <Words className="hello-text">  </Words>
                            <img className="home-unit-image" alt="unit" src={uImage}></img>
                            <Words className="welcome-text">  </Words>
                            <Words animate className="access-text"> </Words>
                            <Words animate className="identification-text"> PLEASE IDENTIFY YOURSELF</Words>
                            <TextField onLogIn={() => { this.props.onLogIn() }}> </TextField>
                        </Content>
                    </Arwes>
                }
            </div>

        );
    }
}

export default withTemplate(withStyles(styles)(LogIn));

