FUNCTION {format.names.author}
{
  's :=
 
  #1 'nameptr :=
  s num.names$ 'numnames :=
 
  numnames 'namesleft :=
  #0 'posothers :=
    { namesleft #0 > }
    { s nameptr "{ll,~}{ff}{vv}{jj}" format.name$ 't :=
      t "others" =
         { nameptr 'posothers := }
         'skip$
       if$
       nameptr #1 + 'nameptr :=
       namesleft #1 - 'namesleft :=
    }
  while$
 
  #0 'showetal :=
 
  numnames #10 >
    { #10 'numnames :=
      #1 'showetal :=
    }
    { numnames #10 =
        { posothers #0 >
          { #1 'showetal :=
            #10 'numnames :=
          }
          'skip$
         if$
         }
         { posothers #0 >
             { #1 'showetal :=
               #1 'numnames -
             }
             'skip$
           if$
         }
       if$
    }
  if$
 
  #1 'nameptr :=
  numnames 'namesleft :=
    { namesleft #0 > }
    { s nameptr "{ll,~}{ff}{vv}{jj}" format.name$ 't :=
      nameptr #1 >
        { namesleft #1 >
            { ", " * t * }
            { showetal #0 >
                { ", " * t * }
                { " and " * t * }
              if$
            }
          if$
        }
        't
      if$
      nameptr #1 + 'nameptr :=
      namesleft #1 - 'namesleft :=
    }
  while$
 
  showetal #0 >
    { " et~al." * }
    'skip$
  if$
}