{{ $personID := .PersonID }}
{{ range $index, $i := .Items }}
  <li><div>
    <div onclick="setSelectedItem('{{$i.ID}}')">{{$i.Title}} ({{countVotes $i.ID }} votes)</div>
    <div>
 {{- if $personID }}
  {{ if haveIUpvoted $i.ID }}
              <img src="https://img.icons8.com/material/24/000000/undo.png" onclick="unvote('{{$i.ID}}')" />
  {{ else }}
              <img src="https://img.icons8.com/material/24/000000/circled-chevron-up.png" onclick="upvote('{{$i.ID}}')" />
  {{ end }}
  {{ if haveIDownvoted $i.ID }}
              <img src="https://img.icons8.com/material/24/000000/undo.png" onclick="unvote('{{$i.ID}}')" />
  {{ else }}
              <img src="https://img.icons8.com/material/24/000000/circled-chevron-down.png" onclick="downvote('{{$i.ID}}')" />
  {{ end }}
 {{ else }}
          <abbr title='log in to vote'><img src="https://img.icons8.com/material/24/000000/question.png" onclick="alert('log in to vote')"></abbr>
 {{ end }}
        </div>
    </div></li>
{{ end }}
