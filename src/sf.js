import { h, patch } from 'https://unpkg.com/superfine?module';
import { countVotes, haveIUpvoted, haveIDownvoted } from './app.js';
import { getPersonId } from './auth-okta.js';

export const view = items => {
  if (items) {
    console.log('items', items);
    const iconPrefix = 'https://img.icons8.com/material/24/000000/';
    return h('ul', { class: 'list' }, [
      items
        .sort((a, b) => countVotes(b.id) - countVotes(a.id))
        .map(i =>
          h('li', {}, [
            h('div', {}, [
              h('div', { onclick: 'setSelectedItem(${i.id})' }, [
                h('div', {}, `${i.title} (${countVotes(i.id)})`),
              ]),
              h(
                'div',
                {},
                getPersonId()
                  ? [
                      h('img', {
                        src:
                          iconPrefix +
                          (haveIUpvoted(i.id)
                            ? 'undo.png'
                            : 'circled-chevron-up.png'),
                        onclick: element => {
                          if (haveIUpvoted(i.id)) {
                            unvote(i.id);
                          } else {
                            upvote('${i.id}');
                          }
                        },
                      }),
                      h('img', {
                        src:
                          iconPrefix +
                          (haveIDownvoted(i.id)
                            ? 'undo.png'
                            : 'circled-chevron-down.png'),
                        onclick: element => {
                          if (haveIDownvoted(i.id)) {
                            unvote(i.id);
                          } else {
                            downvote('${i.id}');
                          }
                        },
                      }),
                    ]
                  : [
                      h(
                        'abbr',
                        { title: 'log in to vote' },
                        h('img', {
                          src: iconPrefix + 'question.png',
                          onclick: element => {
                            alert('log in to vote');
                          },
                        })
                      ),
                    ]
              ),
            ]),
          ])
        ),
    ]);
  } else {
    return h('ul', {}, []);
  }
};
const app = (view, container, node) => items => {
  node = patch(node, view(items), container);
};
let render = null;
export const getRender = vidList => {
  if (render == null) {
    if (vidList != null) {
      render = app(view, vidList);
    }
  }
  return render;
};
