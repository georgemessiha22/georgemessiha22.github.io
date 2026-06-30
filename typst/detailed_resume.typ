#import "@preview/modern-cv:0.10.0": *

#show: resume.with(
  author: (
    firstname: "George",
    lastname: "Messiha",
    email: "georgemessiha22@gmail.com",
    phone: "(+971) 54 555 1032",
    github: "georgemessiha22",
    gitlab: "georgemessiha22",
    linkedin: "georgemessiha22",
    address: "Dubai, UAE",
    positions: ("Lead Software Engineer",),
  ),
  profile-picture: image("assets/profile.jpg"),
  date: datetime.today().display(),
  language: "en",
  colored-headers: true,
  show-footer: false,
  paper-size: "a4",
)

#include "sections/summary.typ"
#include "sections/experience_detailed.typ"
#include "sections/education.typ"
#include "sections/skills.typ"
#include "sections/languages.typ"
#include "sections/certificates_detailed.typ"
#include "sections/activities.typ"
