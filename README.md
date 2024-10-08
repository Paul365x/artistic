# Artistic

Artistic is a file manager for a specific set of use cases. In these cases you manage images with meta data. The images may consist of multiple files. Two immediate use cases are managing artwork for print on demand sites and managing machine embroidery designs. 

Print on Demand artwork consists of a file created and specific to a drawing programme like Adobe Illustrator or the Gimp. To put artwork on a site you need to produce at least one standard graphics file (jpg, png etc) and a range of information such as description, title, tags etc. Managing all this is challenging.

In machine embroidery you have the design file (pes, jef, dst), at least one image file showing the design as a picture and then you may have other files. Again, difficult to manage.

## Personality

Because the use cases have great similarity, there is an opportunity to use mix and match interface elements that meet a configuration for that use case ie personality. This allows specific functionality per use case.

## Design Principles

- You do not need this software to gain access to your artwork, files and metadata. This means you will not lose your work if this software becomes unavailable.

- Metadata will be in a human readable format that is well known in development circles.

- As much as possible, conversion of formats will be easy because the stored format generated by the software will be well understood.

- Reuse is emphasised so that adding use cases is easy.

- Technologies used are cross platform

## Technologies

- Golang

- Fyne UI Framework

- JSON
