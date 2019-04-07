{{ $personID := .PersonID }}
{{ range $index, $i := .Lists }}
<li>
  <div class="video-item" x="{{rand}}" data-id="{{$i.ID}}">
    <div class="title">
      {{$i.Title}} 
    </div>
  </div>
</li>
{{ end }}
