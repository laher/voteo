const hideOkta = () => {
  if (state.oktaSignIn) {
    state.oktaSignIn.hide();
  }
};

const renderOktaSignIn = () => {
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
      state.personId = res[0].claims.email;
      showSignOut();
      getVideos();
      getVotes();
    }
  });
};

const doOkta = () => {
  document.getElementById('sign-out').addEventListener('click', event => {
    event.preventDefault();

    console.log('signout clicked');
    state.oktaSignIn.session.close(err => {
      if (err) {
        alert(`Error: ${err}`);
      }
      showSignIn();
    });
  });
  state.oktaSignIn.session.get(async res => {
    if (res.status === 'ACTIVE') {
      console.log('login already active', res);
      state.personId = res.login;
      state.accessToken = state.oktaSignIn.tokenManager.get(
        'access_token'
      ).accessToken;
      showSignOut();
      getVideos();
      getVotes();
    } else {
      console.log('not signed in');
      showSignIn();
    }
  });
};
