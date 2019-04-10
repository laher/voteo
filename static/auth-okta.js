import { showSignInOut } from './app.js';

let state = {
  personId: null,
  idToken: null,
  accessToken: null,
  oktaSignIn: null,
};

export const hideOkta = () => {
  if (state.oktaSignIn) {
    state.oktaSignIn.hide();
  }
};

export const getAccessToken = () => {
  return state.accessToken;
};

export const getPersonId = () => {
  return state.personId;
};

const renderOktaSignIn = () => {
  document.getElementById('widget-container').style.display = 'none';
  state.oktaSignIn.renderEl({ el: '#widget-container' }, res => {
    if (res.status === 'SUCCESS') {
      console.log('signin success', res);
      state.oktaSignIn.tokenManager.add('id_token', res[0]);
      state.oktaSignIn.tokenManager.add('access_token', res[1]);
      state.idToken = res[0];
      state.accessToken = res[1].accessToken;
      console.log(
        'signin success. tokenManager:',
        state.oktaSignIn.tokenManager
      );
      updateAuthCookie();
      state.personId = res[0].claims.email;
      hideOkta();
      showSignInOut(state.personId);
    }
  });
  state.oktaSignIn.hide();
  document.getElementById('widget-container').style.display = 'block';
};

export const showSignInModal = () => {
  console.log('show signin modal');
  if (state.oktaSignIn) {
    document.getElementById('app-container').style.display = 'none';
    state.oktaSignIn.show();
  } else {
    console.log('oops: non-okta signin not implemented');
  }
};

export const initOkta = (oktaConf, andThen) => {
  state.oktaSignIn = new OktaSignIn(oktaConf);
  document.getElementById('sign-out').addEventListener('click', event => {
    event.preventDefault();
    console.log('signout clicked');
    state.personId = null;
    state.accessToken = null;
    updateAuthCookie();
    state.oktaSignIn.session.close(err => {
      if (err) {
        alert(`Error: ${err}`);
      }
      hideOkta();
      showSignInOut(null);
    });
  });
  renderOktaSignIn();
  state.oktaSignIn.session.get(async res => {
    if (res.status === 'ACTIVE') {
      console.log('login already active', res);
      state.personId = res.login;
      state.accessToken = state.oktaSignIn.tokenManager.get(
        'access_token'
      ).accessToken;
      updateAuthCookie();
      hideOkta();
      showSignInOut(state.personId);
    } else {
      console.log('not signed in');
      hideOkta();
      showSignInOut(null);
    }
    // get (with auth as appropriate - or otherwise)
    andThen();
  });
};

const updateAuthCookie = () => {
  document.cookie = 'auth=' + getAccessToken();
};
